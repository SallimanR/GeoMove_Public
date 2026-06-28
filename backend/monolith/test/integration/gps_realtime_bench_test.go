package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"monolith/cmd/server"
	driverCommands "monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/internal/websockethub"
	wsPB "monolith/internal/websockethub/proto"
	datastructures "monolith/pkg/data_structures"
	"monolith/pkg/logger"
	"monolith/test/testutils"
	"monolith/test/testutils/network"
	testws "monolith/test/testutils/websocket"
)

func setupPublishers(b *testing.B) {}

func BenchmarkGPSRealtimeChannel(b *testing.B) {
	testutils.PingGraphhopper(b)

	logger := logger.NewDebugLogger()

	ctx := context.Background()
	db, dbName, dbCleaunp := testutils.CreateTestDB(b, ctx, adminConn, adminConnString, templateDBName)
	b.Logf("dbName: %s", dbName)

	publisherRole := "tow_driver"
	subscriberRole := "tow_subscriber"
	// wsOptions := websockethub.WebsocketServerOptions{
	// 	Roles: []string{publisherRole, subscriberRole},
	// }
	// wsServer := websockethub.NewWebsocketServer(wsOptions)
	httpPort := network.GetFreePort(b)

	srv := testutils.NewServer(
		b,
		// server.WithWebsocketServer(wsServer),
		server.WithHTTPPort(httpPort),
		server.WithDBConnection(db),
		server.WithDriverDomain(),
		server.WithGeolocationDomain(),
	)
	httpAddrs := srv.GetHttpAddrs()

	dialer := testws.SetupDialer(b, nil)

	testutils.StartServer(b, srv)

	b.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		require.NoError(b, err, "failed to Shutdown server")
		dbCleaunp()
	})

	driverDomain := srv.GetDriverDomain()
	geoDomain := srv.GetGeolocationDomain()
	const channelName = wsPB.Channel_GPS_REALTIME

	const publishersNumber = 100
	const subscribersNumber = 2000

	// Index in starts from 1, not 0
	publishers := make([]*websocket.Conn, publishersNumber+1)
	publishersIndexes := make([]uint32, publishersNumber+1)

	subscribers := make([]*websocket.Conn, subscribersNumber)

	for i := uint32(1); i <= publishersNumber; i++ {
		createDriverCmd := driverCommands.CreateDriverCommand{
			Latitude:  55.5,
			Longitude: 38.4,
		}
		driver, err := driverDomain.Commands.CreateDriver.Handle(ctx, createDriverCmd)
		if err != nil {
			b.Fatalf("failed to create driver: %s", err)
		}

		err = geoDomain.Commands.CreateDriverRealtime.Handle(ctx, command.CreateDriverRealtimeCommand{DriverID: uint32(driver.ID)})
		if err != nil {
			b.Fatalf("failed to create driver in realtime table: %s", err)
		}

		publisherConn := testws.CreateWSConn(b, publisherRole, i, httpAddrs, dialer)

		publishers[i] = publisherConn
		publishersIndexes[i] = uint32(i)
	}

	publisherMessage := pb.LocationUpdate{
		Coordinates: []*pb.Point{
			{Latitude: 55.746055, Longitude: 37.621070},
			{Latitude: 55.745930, Longitude: 37.620733},
			{Latitude: 55.746141, Longitude: 37.620314},
			{Latitude: 55.746582, Longitude: 37.619739},
		},
		Timestamps: []int64{10, 15, 20},
	}
	msg, _ := proto.Marshal(&publisherMessage)
	publisherMsg, _ := proto.Marshal(&wsPB.Request{Channel: channelName, Payload: &wsPB.Request_Publish{
		Publish: &wsPB.PublishMessage{
			Data: msg,
		},
	}})

	type responseMessage struct {
		responseCode uint32
		responseData []*pb.DriverRealtime
	}
	subscribersResponses := datastructures.SyncMap[*websocket.Conn, responseMessage]{}
	for i := range subscribersNumber {
		subscriberConn := testws.CreateWSConn(b, subscriberRole, uint32(i), httpAddrs, dialer)
		subscriberConn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
			respMsg := testws.UnmarshalWSResponseMessage(b, data)
			respData, err := websockethub.UnmarshalProtobufBatch(respMsg.Data, func() *pb.DriverRealtime { return &pb.DriverRealtime{} })
			if err != nil {
				b.Log(err)
			}

			subscribersResponses.Store(conn, responseMessage{responseCode: respMsg.StatusCode, responseData: respData})
		})

		subscribers[i] = subscriberConn
	}

	subscriberMsg, _ := proto.Marshal(&wsPB.Request{Channel: channelName, Payload: &wsPB.Request_Subscribe{
		Subscribe: &wsPB.SubscribeMessage{
			Indexes: publishersIndexes,
		},
	}})

	logger.Info().Msg("Starting benchmark")
	// for b.Loop() {
	b.ResetTimer()
	for i := 1; i <= publishersNumber; i++ {
		err := publishers[i].WriteMessage(websocket.BinaryMessage, publisherMsg)
		if err != nil {
			b.Logf("Failed to write message: %s", err)
		}

	}

	// TODO: publisher connection to bench setup
	b.StopTimer()
	time.Sleep(100 * time.Millisecond)
	b.StartTimer()

	for i := 0; i < subscribersNumber; i++ {
		err := subscribers[i].WriteMessage(websocket.BinaryMessage, subscriberMsg)
		if err != nil {
			b.Logf("Failed to write message: %s", err)
		}
	}
	b.StopTimer()

	time.Sleep(1000 * time.Millisecond)
	var nonfullMessageCount uint32
	var emptyMessageCount uint32
	var badStatusCodeCount uint32
	var totalMessages int
	for i := 0; i < subscribersNumber; i++ {
		resp, ok := subscribersResponses.Load(subscribers[i])
		if !ok {
		}
		messagesCount := len(resp.responseData)
		if messagesCount < publishersNumber {
			nonfullMessageCount += 1
			if messagesCount == 0 {
				emptyMessageCount += 1
			}
		}
		if resp.responseCode != http.StatusOK {
			badStatusCodeCount += 1
		}
		totalMessages += messagesCount
	}
	resp, ok := subscribersResponses.Load(subscribers[subscribersNumber-1])
	if !ok {
	}
	logger.Debug().
		Uint32("response code", resp.responseCode).
		Int("response data lenght", len(resp.responseData)).
		Any("response data", resp.responseData).
		Send()

	logger.Debug().
		Uint32("nonfullMessageCount", nonfullMessageCount).
		Uint32("emptyMessageCount", emptyMessageCount).
		Uint32("badStatusCodeCount", badStatusCodeCount).
		Int("totalMessages", totalMessages).
		Float64("averageMessages", float64(totalMessages)/float64(subscribersNumber)).
		Send()

	// TODO: cleanup connections:
	// b.Cleanup(func() {
	// 	for i := 1; i <= publishersNumber; i++ {
	// 		publisherConn := publishers[i]
	// 		err = publisherConn.Close()
	// 		if err != nil {
	// 			log.Println("Failed to close client connection")
	// 		}
	// 	}
	// })
	//
	// err = subscriberConn.Close()
	// if err != nil {
	// 	log.Println("[test] failed to close client connection")
	// }
}
