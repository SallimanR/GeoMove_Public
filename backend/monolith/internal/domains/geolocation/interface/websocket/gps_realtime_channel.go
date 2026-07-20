package websocket

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/infrastructure/graphopper"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/internal/websockethub"
	wsPB "monolith/internal/websockethub/proto"
)

type GPSRealtimeChannel struct {
	Channel *websockethub.PubSubChannel[entity.MovingDriverWithPoints]

	updateMovingDriver *command.UpdateMovingDriverHandler

	logger zerolog.Logger
}

func NewGPSRealtimeChannel(wsServer *websockethub.WebsocketServer, roles []string, updateMovingDriver *command.UpdateMovingDriverHandler, logger zerolog.Logger) (*GPSRealtimeChannel, wsPB.Channel, error) {
	const channelName = wsPB.Channel_GPS_REALTIME
	channel := &GPSRealtimeChannel{
		Channel:            websockethub.NewPubSubChannel[entity.MovingDriverWithPoints](),
		updateMovingDriver: updateMovingDriver,
		logger:             logger,
	}
	err := wsServer.RegisterChannel(roles, channelName, channel)
	if err != nil {
		return nil, channelName, err
	}
	return channel, channelName, nil
}

func (c *GPSRealtimeChannel) Publish(publisherID int64, msg []byte) error {
	data := pb.GPSUpdate{}
	err := proto.Unmarshal(msg, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message")
	}

	matchReq, err := graphopper.BuildGraphhopperRequestFromWS(data.Coordinates, data.Timestamps)
	if err != nil {
		return err
	}
	c.logger.Trace().Any("matchReq", matchReq).Send()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	matchResp, ok := graphopper.MatchMapCoordinates(ctx, matchReq)
	if !ok {
		return fmt.Errorf("failed to match coordinates")
	}
	c.logger.Trace().Any("matchResp", matchResp).Send()

	coordinatesNumber := len(matchResp.Points.Coordinates)
	if coordinatesNumber < 2 {
		return nil
	}

	lastCoord := matchResp.Points.Coordinates[coordinatesNumber-1]
	lon := lastCoord[0]
	lat := lastCoord[1]

	points := make([][2]float32, coordinatesNumber)
	for i, c := range matchResp.Points.Coordinates {
		points[i] = [2]float32{c[0], c[1]}
	}

	travelDuration := time.Duration(matchResp.Time) * time.Millisecond

	gpsData := entity.MovingDriverWithPoints{
		MovingDriver: entity.MovingDriver{
			DriverID:   publisherID,
			Latitude:   lat,
			Longitude:  lon,
			TravelTime: time.Time{}.Add(travelDuration),
			PathMeters: int32(matchResp.Distance),
		},
		Points: points,
	}

	c.Channel.Messages.Store(publisherID, gpsData)

	err = c.updateMovingDriver.Handle(ctx, command.UpdateMovingDriverCommand{
		DriverID:    publisherID,
		Coordinates: matchResp.Points.Coordinates,
		TravelTime:  time.Time{}.Add(travelDuration),
		PathMeters:  uint32(matchResp.Distance),
	})
	if err != nil {
		c.logger.Debug().AnErr("Channel_GPS_REALTIME err", err).Send()
		return fmt.Errorf("failed to update location in DB: %s", err)
	}

	return nil
}

func (c *GPSRealtimeChannel) GetMessages(publisherIDs []int64) ([]byte, error) {
	items := c.Channel.GetMessages(publisherIDs)

	messages := make([][]byte, 0, len(items))
	for _, item := range items {
		points := make([]*pb.Coordinates, len(item.Points))
		for i, p := range item.Points {
			points[i] = &pb.Coordinates{
				Coordinates: []float32{p[0], p[1]},
			}
		}

		md := &pb.MovingDriver{
			DriverId:   item.DriverID,
			Lat:        item.Latitude,
			Lon:        item.Longitude,
			TravelTime: uint32(item.TravelTime.Sub(time.Time{}).Seconds()),
			PathMeters: uint32(item.PathMeters),
			Points:     points,
		}
		msgProto, err := proto.Marshal(md)
		if err != nil {
			continue
		}
		messages = append(messages, msgProto)
	}

	var messagesProtobuf wsPB.MessageBatch
	messagesProtobuf.Data = messages
	return proto.Marshal(&messagesProtobuf)
}
