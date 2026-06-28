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
	"monolith/internal/domains/geolocation/application/query"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/internal/websockethub"
	wsPB "monolith/internal/websockethub/proto"
	"monolith/pkg/logger"
	"monolith/test/testutils"
	testDomains "monolith/test/testutils/domains"
	"monolith/test/testutils/network"
	testws "monolith/test/testutils/websocket"
)

// TODO: token auth
// func createWSConnection(t *testing.T, tocken string, wsAddr string, dialer *websocket.Dialer) *websocket.Conn {
func createWSConn(t *testing.T, role string, id uint32, httpAddr string, dialer *websocket.Dialer) (*websocket.Conn, chan []byte) {
	t.Helper()
	conn := testws.CreateWSConn(t, role, id, httpAddr, dialer)
	response := make(chan []byte, 1)
	conn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
		response <- data
	})
	return conn, response
}

func TestGPSRaltimeChannel_Integration(t *testing.T) {
	testutils.PingGraphhopper(t)

	logger := logger.NewDebugLogger()

	ctx := context.Background()
	db, dbName, dbCleaunp := testutils.CreateTestDB(t, ctx, adminConn, adminConnString, templateDBName)
	t.Logf("dbName: %s", dbName)

	publisherRole := "tow_driver"
	subscriberRole := "tow_subscriber"
	// wsOptions := websockethub.WebsocketServerOptions{
	// 	Roles: []string{publisherRole, subscriberRole},
	// }
	// wsServer := websockethub.NewWebsocketServer(wsOptions)
	httpPort := network.GetFreePort(t)

	srv := testutils.NewServer(
		t,
		// server.WithWebsocketServer(wsServer),
		server.WithHTTPPort(httpPort),
		server.WithDBConnection(db),
		server.WithDriverDomain(),
		server.WithGeolocationDomain(),
	)
	httpAddrs := srv.GetHttpAddrs()

	dialer := testws.SetupDialer(t, nil)

	testutils.StartServer(t, srv)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		require.NoError(t, err, "failed to Shutdown server")
		dbCleaunp()
	})

	driverDomain := srv.GetDriverDomain()
	createDriverCmd := driverCommands.CreateDriverCommand{
		Latitude:  55.5,
		Longitude: 38.4,
	}
	driver := testDomains.CreateDriver(t, driverDomain.Commands.CreateDriver, createDriverCmd)
	geoDomain := srv.GetGeolocationDomain()
	const channelName = wsPB.Channel_GPS_REALTIME
	// const channelNameFake wsPB.Channel = 1
	testDomains.CreateDriverRealtime(t, geoDomain, command.CreateDriverRealtimeCommand{DriverID: uint32(driver.ID)})

	testCases := []struct {
		name string

		publisherID      uint32
		publisherRole    string
		publisherData    pb.LocationUpdate
		publisherWantErr bool

		subscriberID      uint32
		subscriberRole    string
		subscriberData    []uint32
		subscriberWantErr bool
	}{
		{
			name: "publisher/subscriber id, role, data is correct",

			// TODO: database seeding
			publisherID:   uint32(driver.ID),
			publisherRole: publisherRole,
			publisherData: pb.LocationUpdate{
				Coordinates: []*pb.Point{
					{Latitude: 55.746055, Longitude: 37.621070},
					{Latitude: 55.745930, Longitude: 37.620733},
					{Latitude: 55.746141, Longitude: 37.620314},
					{Latitude: 55.746582, Longitude: 37.619739},
				},
				Timestamps: []int64{10, 15, 20},
			},

			// TODO: unregistered subscribets would don't have ids, and therefore role too
			subscriberID:   10,
			subscriberRole: subscriberRole,
			subscriberData: []uint32{uint32(driver.ID)},
		},
	}
	for i := 0; i < len(testCases) && !t.Failed(); i++ {
		tc := &testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			publisherConn, publisherResponse := createWSConn(t, tc.publisherRole, tc.publisherID, httpAddrs, dialer)

			msg, _ := proto.Marshal(&tc.publisherData)
			publisherMsg, _ := proto.Marshal(&wsPB.Request{Channel: channelName, Payload: &wsPB.Request_Publish{
				Publish: &wsPB.PublishMessage{
					Data: msg,
				},
			}})
			testws.WriteMessage(t, publisherConn, publisherMsg)

			select {
			case data := <-publisherResponse:
				respMsg := testws.UnmarshalWSResponseMessage(t, data)
				logger.Debug().Any("[test] publisher response", &respMsg).Send()
				if !tc.publisherWantErr {
					require.EqualValues(t, http.StatusOK, respMsg.StatusCode, "publisher send wrong message")
				}
			case <-time.After(1000 * time.Millisecond):
				t.Fatalf("timeout writing message")
			}

			_, err := geoDomain.Queries.GetDriverByID.Handle(ctx, query.GetDriverRealtimeByIDQuery{DriverID: tc.publisherID})
			require.NoErrorf(t, err, "failed to find realtime by id")

			subscriberConn, subscriberResponse := createWSConn(t, tc.subscriberRole, tc.subscriberID, httpAddrs, dialer)
			subscriberMsg, _ := proto.Marshal(&wsPB.Request{Channel: channelName, Payload: &wsPB.Request_Subscribe{
				Subscribe: &wsPB.SubscribeMessage{
					Indexes: tc.subscriberData,
				},
			}})

			testws.WriteMessage(t, subscriberConn, subscriberMsg)
			respTime := time.Now()
			for {
				select {
				case data := <-subscriberResponse:
					respTimeDelta := time.Since(respTime)
					logger.Debug().Int64("respTimeDelta", respTimeDelta.Milliseconds()).Send()
					respTime = time.Now()

					respMsg := testws.UnmarshalWSResponseMessage(t, data)

					responseData, err := websockethub.UnmarshalProtobufBatch(respMsg.Data, func() *pb.DriverRealtime { return &pb.DriverRealtime{} })
					require.NoError(t, err, "failed to unmarshal consumer response data")
					logger.Debug().Uint32("[test] consumer response status", respMsg.StatusCode).Any("[test] consumer response data", responseData).Send()
					if !tc.subscriberWantErr {
						require.EqualValues(t, http.StatusOK, respMsg.StatusCode, "subscriber send wrong message")
						require.Len(t, responseData, 1, "subscriber did not received any data")
					}
					// case <-time.After(1000 * time.Millisecond):
					// 	t.Fatalf("timeout writing message")
				}
			}
		})
	}
}
