package testws

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"

	"monolith/test/testutils/network"
)

type WSDialer struct {
	*websocket.Dialer
}

func SetupDialer(t testing.TB, opts *websocket.Upgrader) *websocket.Dialer {
	dialerPort := network.GetFreePort(t)
	dialerAddr := fmt.Sprintf("127.0.0.1:%d", dialerPort)
	httpMux := &http.ServeMux{}
	wsEngine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{dialerAddr},
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 httpMux,
	})
	err := wsEngine.Start()
	if err != nil {
		t.Fatalf("failed to start engine: %s", err)
	}
	if opts == nil {
		opts = &websocket.Upgrader{}
	}
	dialer := &websocket.Dialer{
		Engine:            wsEngine,
		DialTimeout:       5 * time.Second,
		EnableCompression: false,
		Upgrader:          opts,
	}
	return dialer
}

func ShutdownDialer(dialer *websocket.Dialer) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := dialer.Engine.Shutdown(ctx)
	if err != nil {
	}
}
