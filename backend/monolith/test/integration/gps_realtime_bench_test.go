package integration

import (
	"context"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"monolith/cmd/server"
	driverCommands "monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/internal/notification"
	"monolith/internal/websockethub"
	wsPB "monolith/internal/websockethub/proto"
	datastructures "monolith/pkg/data_structures"
	"monolith/pkg/logger"
	"monolith/test/testutils"
	testauth "monolith/test/testutils/auth"
	"monolith/test/testutils/network"
	testws "monolith/test/testutils/websocket"
)

func setupPublishers(b *testing.B) {}

func BenchmarkGPSRealtimeChannel(b *testing.B) {
	testutils.PingGraphhopper(b)

	vapidKeys, err := notification.GenerateVAPIDKeys()
	require.NoError(b, err)
	os.Setenv("VAPID_PUBLIC_KEY", vapidKeys.PublicKey)
	os.Setenv("VAPID_PRIVATE_KEY", vapidKeys.PrivateKey)

	logger := logger.NewDebugLogger()

	ctx := context.Background()
	db, dbName, dbCleaunp := testutils.CreateTestDB(b, ctx, adminConn, adminConnString, templateDBName)
	b.Logf("dbName: %s", dbName)

	publisherRole := "tow_driver"
	subscriberRole := "tow_subscriber"

	httpPort := network.GetFreePort(b)

	srv := testutils.NewServer(
		b,
		server.WithHTTPPort(httpPort),
		server.WithDBConnection(db),
		server.WithAuth(),
		server.WithNotification(),
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

	const channelName = wsPB.Channel_GPS_REALTIME

	const publishersNumber = 100
	const subscribersNumber = 2000

	publisherUsers := testauth.CreateUsersWithSessionBatch(ctx, b, db, []string{publisherRole, subscriberRole}, publishersNumber)
	subscriberUsers := testauth.CreateUsersWithSessionBatch(ctx, b, db, []string{subscriberRole}, subscribersNumber)

	var wg sync.WaitGroup
	var driverErr atomic.Int64
	for _, user := range publisherUsers {
		wg.Add(1)
		go func(user *testauth.TestUser) {
			defer wg.Done()
			createDriverCmd := driverCommands.CreateDriverCommand{
				UserID:    user.ID,
				Name:      "bench-driver",
				Latitude:  55.5,
				Longitude: 38.4,
			}
			err := driverDomain.Commands.CreateDriver.Handle(ctx, createDriverCmd)
			if err != nil {
				driverErr.Add(1)
				b.Logf("failed to create driver: %s", err)
			}
		}(user)
	}
	wg.Wait()
	if n := driverErr.Load(); n > 0 {
		b.Fatalf("%d drivers failed to create", n)
	}

	publishers := make([]*websocket.Conn, publishersNumber)
	publishersIndexes := make([]int64, publishersNumber)
	publisherDone := make(chan struct{}, publishersNumber)
	var publisherErrCount atomic.Int64

	for i := range publishersNumber {
		wsHeaders := http.Header{}
		wsHeaders.Set("Cookie", "session="+publisherUsers[i].SessionToken)
		publisherConn := testws.CreateWSConn(b, publisherRole, publisherUsers[i].ID, httpAddrs, dialer, wsHeaders)
		publisherConn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
			var resp wsPB.Response
			if err := proto.Unmarshal(data, &resp); err != nil {
				publisherErrCount.Add(1)
				publisherDone <- struct{}{}
				return
			}
			if respMsg, ok := resp.Payload.(*wsPB.Response_Response); ok {
				if respMsg.Response.StatusCode != http.StatusOK {
					publisherErrCount.Add(1)
				}
			} else {
				publisherErrCount.Add(1)
			}
			publisherDone <- struct{}{}
		})
		publishers[i] = publisherConn
		publishersIndexes[i] = publisherUsers[i].ID
	}

	publisherMessage := pb.GPSUpdate{
		Coordinates: []*pb.Point{
			{Latitude: 55.746055, Longitude: 37.621070},
			{Latitude: 55.745930, Longitude: 37.620733},
			{Latitude: 55.746141, Longitude: 37.620314},
			{Latitude: 55.746582, Longitude: 37.619739},
		},
		Timestamps: []int64{10, 15, 20},
	}
	msg, _ := proto.Marshal(&publisherMessage)
	publisherMsg, _ := proto.Marshal(&wsPB.Request{RequestId: "pub", Channel: channelName, Payload: &wsPB.Request_Publish{
		Publish: &wsPB.PublishMessage{
			Data: msg,
		},
	}})

	logger.Info().Msg("Publishing GPS data")
	for i := range publishersNumber {
		err := publishers[i].WriteMessage(websocket.BinaryMessage, publisherMsg)
		if err != nil {
			b.Logf("Failed to write message: %s", err)
		}
	}

	for range publishersNumber {
		select {
		case <-publisherDone:
		case <-time.After(30 * time.Second):
			b.Fatal("timeout waiting for publisher responses")
		}
	}

	if n := publisherErrCount.Load(); n > 0 {
		b.Logf("%d/%d publishers had errors", n, publishersNumber)
	}

	type responseMessage struct {
		responseCode uint32
		responseData []*pb.MovingDriver
	}
	subscribersResponses := datastructures.SyncMap[*websocket.Conn, responseMessage]{}
	subscribers := make([]*websocket.Conn, subscribersNumber)

	for i := range subscribersNumber {
		subHeaders := http.Header{}
		subHeaders.Set("Cookie", "session="+subscriberUsers[i].SessionToken)
		subscriberConn := testws.CreateWSConn(b, subscriberRole, subscriberUsers[i].ID, httpAddrs, dialer, subHeaders)
		subscriberConn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
			var resp wsPB.Response
			if err := proto.Unmarshal(data, &resp); err != nil {
				return
			}
			respMsg, ok := resp.Payload.(*wsPB.Response_Response)
			if !ok {
				return
			}
			if respMsg.Response.StatusCode != http.StatusOK {
				return
			}
			respData, err := websockethub.UnmarshalProtobufBatch(respMsg.Response.Data, func() *pb.MovingDriver { return &pb.MovingDriver{} })
			if err != nil || len(respData) == 0 {
				return
			}
			subscribersResponses.Store(conn, responseMessage{responseCode: respMsg.Response.StatusCode, responseData: respData})
		})
		subscribers[i] = subscriberConn
	}

	subscriberMsg, _ := proto.Marshal(&wsPB.Request{RequestId: "sub", Channel: channelName, Payload: &wsPB.Request_Subscribe{
		Subscribe: &wsPB.SubscribeMessage{
			Indexes: publishersIndexes,
		},
	}})

	b.ResetTimer()
	for i := range subscribersNumber {
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
	for i := range subscribersNumber {
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
		Int("response data length", len(resp.responseData)).
		Any("response data", resp.responseData).
		Send()

	logger.Debug().
		Uint32("nonfullMessageCount", nonfullMessageCount).
		Uint32("emptyMessageCount", emptyMessageCount).
		Uint32("badStatusCodeCount", badStatusCodeCount).
		Int("totalMessages", totalMessages).
		Float64("averageMessages", float64(totalMessages)/float64(subscribersNumber)).
		Send()
}
