package websocket

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/infrastructure/graphopper"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/internal/websockethub"
	wsPB "monolith/internal/websockethub/proto"
)

type GPSRealtimeChannel struct {
	Channel  *websockethub.PubSubChannel
	wsServer *websockethub.WebsocketServer

	updateDriverRealtime *command.UpdateDriverRealtimeHandler

	// TODO: queues and batching (and benchmark: 1. without batching => go routine per req. 2 with batching)
	// matchCh          chan matchWork
	//    dbCh             chan dbWork
	//    stopCh           chan struct{}
	//    wg               sync.WaitGroup

	logger    zerolog.Logger
	debugMode bool
}

// TODO: use struct instead of parameters?
func NewGPSRealtimeChannel(wsServer *websockethub.WebsocketServer, roles []string, updateDriverRealtime *command.UpdateDriverRealtimeHandler, logger zerolog.Logger) (*GPSRealtimeChannel, wsPB.Channel, error) {
	const channelName = wsPB.Channel_GPS_REALTIME
	channel := &GPSRealtimeChannel{
		Channel:              websockethub.NewPubSubChannel(),
		updateDriverRealtime: updateDriverRealtime,
		logger:               logger,
	}
	err := wsServer.RegisterChannel(roles, channelName, channel)
	if err != nil {
		return nil, channelName, err
	}
	return channel, channelName, nil
}

// TODO: createDriverRealtime command
func Subscribe() {}

func (c *GPSRealtimeChannel) Publish(publisherID uint32, msg []byte) error {
	data := pb.LocationUpdate{}
	err := proto.Unmarshal(msg, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message")
	}

	matchReq, err := graphopper.BuildGraphhopperRequestFromWS(data.Coordinates, data.Timestamps)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	matched, ok := graphopper.MatchMapCoordinates(ctx, matchReq)
	if !ok {
		return fmt.Errorf("failed to match coordinates")
	}

	coordinatesNumber := len(matched.Points.Coordinates)
	// TODO: 1. implement staying mode for car
	// (car could stay on red light, then we may don't update location, but change state "staying"/"driving")
	//       2. move validation to domain entities?
	//
	// var isStaying bool
	if coordinatesNumber < 2 {
		return nil
	}

	// msgPB := ports.DriverRealtime{
	// 	DriverId: publisherID,
	// 	Distance: matched.Distance,
	// 	Time:     uint32(matched.Time),
	// 	Points:   make([]float32, coordinatesNumber*2),
	// }
	// for i := 0; i < len(matched.Points.Coordinates); i++ {
	// 	msgPB.Points[i*2] = matched.Points.Coordinates[i][1]
	// 	msgPB.Points[i*2+1] = matched.Points.Coordinates[i][0]
	// }

	msgPB := pb.DriverRealtime{
		DriverId: publisherID,
		Distance: matched.Distance,
		Time:     uint32(matched.Time),
		Points:   make([]*pb.Coordinates, coordinatesNumber),
	}
	for i, coordinates := range matched.Points.Coordinates {
		msgPB.Points[i] = &pb.Coordinates{Coordinates: coordinates[:]}
	}
	msgMarshalled, err := proto.Marshal(&msgPB)
	if err != nil {
		return fmt.Errorf("failed to marshal: %s", err)
	}
	c.Channel.Messages.Store(publisherID, msgMarshalled)

	err = c.updateDriverRealtime.Handle(ctx, command.UpdateDriverRealtimeCommand{
		DriverID:    publisherID,
		Coordinates: matched.Points.Coordinates,
		Time:        uint64(matched.Time),
		Distance:    matched.Distance,
	})
	if err != nil {
		return fmt.Errorf("failed to update location in DB: %s", err)
	}

	return nil
}

// func (lc *LocationChannel) handleLocationUpdate(publisherID uint32, msg []byte) {
// }

func (c *GPSRealtimeChannel) GetMessages(publisherIDs []uint32) ([]byte, error) {
	messages, err := c.Channel.GetMessages(publisherIDs)
	return messages, err
}
