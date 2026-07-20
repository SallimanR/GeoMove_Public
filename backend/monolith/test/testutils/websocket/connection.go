package testws

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func CreateWSConn(t testing.TB, role string, id int64, httpAddr string, dialer *websocket.Dialer, headers ...http.Header) *websocket.Conn {
	t.Helper()
	connURL := fmt.Sprintf("ws://%s/ws/%s?id=%d", httpAddr, role, id)
	h := http.Header{}
	if len(headers) > 0 {
		h = headers[0]
	}
	conn, res, err := dialer.Dial(connURL, h)
	if err != nil {
		if res != nil && res.Body != nil {
			bReason, _ := io.ReadAll(res.Body)
			t.Fatalf("dial failed: %v, reason: %v\n", err, string(bReason))
		}
		t.Fatalf("Client %d: Failed to connect: %v", id, err)
	}

	return conn
}
