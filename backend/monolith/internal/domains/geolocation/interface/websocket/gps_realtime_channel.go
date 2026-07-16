package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/infrastructure/graphopper"
	geoHTTP "monolith/internal/domains/geolocation/interface/http"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/internal/websockethub"
	wsPB "monolith/internal/websockethub/proto"
)

type GPSRealtimeChannel struct {
	Channel *websockethub.PubSubChannel

	updateMovingDriver *command.UpdateMovingDriverHandler

	logger zerolog.Logger
}

func NewGPSRealtimeChannel(wsServer *websockethub.WebsocketServer, roles []string, updateMovingDriver *command.UpdateMovingDriverHandler, logger zerolog.Logger) (*GPSRealtimeChannel, wsPB.Channel, error) {
	const channelName = wsPB.Channel_GPS_REALTIME
	channel := &GPSRealtimeChannel{
		Channel:            websockethub.NewPubSubChannel(),
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

	points := make([][]float32, coordinatesNumber)
	for i, c := range matchResp.Points.Coordinates {
		points[i] = []float32{c[0], c[1]}
	}

	dr := geoHTTP.MovingDriver{
		DriverId:   int(publisherID),
		Lat:        lat,
		Lon:        lon,
		TravelTime: float32(matchResp.Time),
		PathMeters: int(matchResp.Distance),
		Points:     points,
	}

	msgJSON, err := json.Marshal(dr)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %s", err)
	}
	c.Channel.Messages.Store(publisherID, msgJSON)

	travelDuration := time.Duration(matchResp.Time) * time.Millisecond
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
	messages, err := c.Channel.GetMessages(publisherIDs)
	return messages, err
}
