package testutils

import (
	"os"
	"testing"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/zerolog"

	"monolith/cmd/server"
	"monolith/internal/websockethub"
	"monolith/test/testutils/network"
)

type WSServer struct {
	server   *server.Server
	wsServer *websockethub.WebsocketServer
	dialer   *websocket.Dialer

	httpAddrs string
	logger    zerolog.Logger
}

func NewWSServer(t *testing.T, wsOptions websockethub.WebsocketServerOptions) *WSServer {
	t.Helper()
	wsServer := websockethub.NewWebsocketServer(wsOptions)

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()

	wsPort := network.GetFreePort(t)
	srv, err := server.NewServer(
		server.WithWebsocketServer(wsServer),
		server.WithHTTPPort(wsPort),
	)
	if err != nil {
		t.Fatalf("Failed to setup server: %s", err)
	}
	httpAddrs := srv.GetHttpAddrs()
	logger.Debug().Str("httpAddrs", httpAddrs).Send()

	upgrader := &websocket.Upgrader{}
	dialer := &websocket.Dialer{
		// Engine:            srv.GetWSEngine(),
		DialTimeout:       5 * time.Second,
		EnableCompression: false,
		Upgrader:          upgrader,
	}

	return &WSServer{
		server:    srv,
		wsServer:  wsServer,
		dialer:    dialer,
		httpAddrs: httpAddrs,
		logger:    logger,
	}
}
