package websockethub

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"monolith/test/testutils/network"
	testws "monolith/test/testutils/websocket"
)

type mockSession struct {
	userID int64
	roles  []string
}

func (m *mockSession) GetUserID() int64   { return m.userID }
func (m *mockSession) GetRoles() []string { return m.roles }

func startHTTPRouter(t testing.TB, wsServer *WebsocketServer, defaultRole string) string {
	t.Helper()
	httpRouter := gin.New()
	httpRouter.Use(func(c *gin.Context) {
		role := c.Param("role")
		userID := int64(1)
		if idStr := c.Query("id"); idStr != "" {
			if parsed, err := fmt.Sscanf(idStr, "%d", &userID); err != nil || parsed != 1 {
				userID = 1
			}
		}
		mock := &mockSession{
			userID: userID,
			roles:  []string{role, defaultRole},
		}
		c.Set("user", mock)
	})
	httpRouter.GET("/ws/:role", wsServer.WebsocketUpgradeHandler)
	wsPort := network.GetFreePort(t)
	wsAddr := fmt.Sprintf("127.0.0.1:%d", wsPort)
	go func() {
		err := httpRouter.Run(wsAddr)
		if err != nil {
			t.Fail()
		}
	}()
	time.Sleep(100 * time.Millisecond)

	return wsAddr
}

func TestWebsocketConnectionValidity(t *testing.T) {
	t.Parallel()
	// TODO: authorization tests
	role := "test"

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	wsOptions := WebsocketServerOptions{
		Roles:  []string{"test"},
		Logger: logger,
	}
	wsServer := NewWebsocketServer(wsOptions)
	dialer := testws.SetupDialer(t, nil)
	defer testws.SetupDialer(t, nil)
	wsAddr := startHTTPRouter(t, wsServer, role)

	testCases := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "RoleExists",
			role:    role,
			wantErr: false,
		},
		{
			name:    "RoleDoesNotExists",
			role:    "penguin",
			wantErr: true,
		},
	}
	for i := 0; i < len(testCases) && !t.Failed(); i++ {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			conn, res, err := dialer.Dial(fmt.Sprintf("ws://%s/ws/%s", wsAddr, tc.role), http.Header{})
			if tc.wantErr {
				require.Error(t, err, "Upgraded non existing role")
				return
			}
			if err != nil {
				if res != nil && res.Body != nil {
					bReason, _ := io.ReadAll(res.Body)
					t.Fatalf("dial failed: %v, reason: %v\n", err, string(bReason))
				}
				t.Fatalf("Client %d: Failed to connect: %v", 1, err)
			}

			time.Sleep(1000 * time.Microsecond)
			connPool, _ := wsServer.ConnectionsByRole[tc.role]
			_, ok := connPool.activeConnections.Load(1)
			require.True(t, ok, "No connection in activeConnections")

			err = conn.Close()
			require.NoError(t, err, "Failed to close client connection")

			time.Sleep(10 * time.Millisecond)
			_, ok = connPool.activeConnections.Load(1)
			require.False(t, ok, "Connection in activeConnections, after it was closed")
		})
	}
}

func TestWebsocketConnectionPersistence(t *testing.T) {
	t.Parallel()
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	role := "test"
	wsOptions := WebsocketServerOptions{
		Roles:  []string{role},
		Logger: logger,
	}
	wsServer := NewWebsocketServer(wsOptions)
	dialer := testws.SetupDialer(t, nil)
	wsAddr := startHTTPRouter(t, wsServer, role)

	// WARNING: increasing number of conenction causing timeouts, because lastency increases
	connNumber := 10

	conns := make([]*websocket.Conn, connNumber)
	var timeoutCount uint32

	// TODO: workers instead of go routine per connection
	for i := 0; i < connNumber; i++ {
		conn := testws.CreateWSConn(t, role, uint32(i), wsAddr, dialer)
		conns[i] = conn

		done := make(chan bool)
		conn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
			close(done)
		})

		err := conn.WriteMessage(websocket.BinaryMessage, []byte("test connection"))
		require.NoError(t, err, fmt.Sprintf("Failed to write message, userID: %d", i))

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			timeoutCount++
		}
	}

	log.Printf("timeout count: %d", timeoutCount)
	connPool := wsServer.ConnectionsByRole[role]

	for i := uint32(0); i < uint32(connNumber); i++ {
		_, ok := connPool.activeConnections.Load(i)
		require.True(t, ok, "No connection in activeConnections")
	}

	for _, conn := range conns {
		_ = conn.Close()
	}
	// TODO: proper synchronization
	time.Sleep(100 * time.Millisecond)
	for i := uint32(0); i < uint32(connNumber); i++ {
		_, ok := connPool.activeConnections.Load(i)
		require.False(t, ok, "connection is in activeConnections after it was closed")
	}
}

// func TestWebsocketMessage(t *testing.T) {
// 	t.Parallel()
// 	output := zerolog.ConsoleWriter{Out: os.Stdout}
// 	logger := zerolog.New(output).
// 		Level(zerolog.InfoLevel).
// 		With().
// 		Timestamp().
// 		Logger()
//
// 	role := "test"
// 	wsOptions := WebsocketServerOptions{
// 		Roles:  []string{role},
// 		Logger: logger,
// 	}
// 	wsServer := NewWebsocketServer(wsOptions)
// 	dialer := testws.SetupDialer(t)
// 	defer testws.SetupDialer(t)
// 	wsAddr := startHTTPRouter(t, wsServer, role)
//
// 	testCases := []struct {
// 		name     string
// 		userID   uint32
// 		userRole string
// 		message  wsPB.Request
// 	}{
// 		{
// 			name:     "Message Test",
// 			userID:   1,
// 			userRole: "test",
// 			message: wsPB.Request{
// 				Channel: 0,
// 				Payload: &wsPB.Request_Publish{
// 					Publish: &wsPB.PublishMessage{
// 						Data: []byte("hello"),
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for i := 0; i < len(testCases) && !t.Failed(); i++ {
// 		tc := &testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			conn := testws.CreateWSConn(t, role, uint32(i), wsAddr, dialer)
// 			gotResponse := make(chan []byte, 1)
// 			conn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
// 				gotResponse <- data
// 			})
//
// 			msg, err := proto.Marshal(&tc.message)
// 			require.NoError(t, err, "Failed to marshal message")
// 			err = conn.WriteMessage(websocket.BinaryMessage, msg)
// 			require.NoError(t, err, "Failed to write message")
//
// 			select {
// 			case resp := <-gotResponse:
// 				var respMsg wsPB.Response
// 				err := proto.Unmarshal(resp, &respMsg)
// 				require.NoError(t, err, "Failed to unmarshal message")
// 				logger.Debug().Any("response message", &respMsg).Send()
// 			case <-time.After(100 * time.Millisecond):
// 				t.Fatalf("timeout writing message")
// 			}
// 		})
// 	}
// }

func BenchmarkWebsocketConnection(b *testing.B) {
	// TODO: see nats/jetstream_benchmark_test.go: 252
	// benchmarkCases := []struct {
	// 	connectionsNumber int
	// 	messageNumber     int
	// 	messageSize       int
	// }{
	// 	{},
	// }
	role := "test"

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	wsOptions := WebsocketServerOptions{
		Roles:  []string{"test"},
		Logger: logger,
	}
	wsServer := NewWebsocketServer(wsOptions)
	dialer := testws.SetupDialer(b, nil)
	wsAddr := startHTTPRouter(b, wsServer, role)

	connNumber := 10000

	messageLatency := make([]time.Duration, connNumber)
	var timeoutCount atomic.Int32
	// TODO: tracing and profiling
	// trace.NewFlightRecorder
	b.ResetTimer()
	for i := 0; i < connNumber; i++ {
		conn := testws.CreateWSConn(b, role, uint32(i), wsAddr, dialer)

		done := make(chan bool)
		conn.OnMessage(func(conn *websocket.Conn, msgType websocket.MessageType, data []byte) {
			close(done)
		})

		err := conn.WriteMessage(websocket.TextMessage, []byte("ping message"))
		if err != nil {
			log.Println("Failed to write message")
		}
		timeBefore := time.Now()

		select {
		case <-done:
			messageLatency[i] = time.Duration(time.Since(timeBefore).Nanoseconds())
		case <-time.After(100 * time.Millisecond):
			timeoutCount.Add(1)
		}
	}
	b.StopTimer()

	sort.Slice(messageLatency, func(i int, j int) bool { return messageLatency[i] < messageLatency[j] })
	// TODO: add average latency
	// NOTE: Latency is in nanoseconds
	latencyP50 := messageLatency[int(float64(len(messageLatency))*0.50)]
	latencyP75 := messageLatency[int(float64(len(messageLatency))*0.75)]
	latencyP90 := messageLatency[int(float64(len(messageLatency))*0.90)]
	latencyP95 := messageLatency[int(float64(len(messageLatency))*0.95)]
	latencyP99 := messageLatency[int(float64(len(messageLatency))*0.99)]
	latencyP99_9 := messageLatency[int(float64(len(messageLatency))*0.999)]

	b.ReportMetric(float64(latencyP50)/float64(time.Microsecond), "p50")
	b.ReportMetric(float64(latencyP75)/float64(time.Microsecond), "p75")
	b.ReportMetric(float64(latencyP90)/float64(time.Microsecond), "p90")
	b.ReportMetric(float64(latencyP95)/float64(time.Microsecond), "p95")
	b.ReportMetric(float64(latencyP99)/float64(time.Microsecond), "p99")
	b.ReportMetric(float64(latencyP99_9)/float64(time.Microsecond), "p99.9")

	log.Println("connections timeout count: ", timeoutCount.Load())
}
