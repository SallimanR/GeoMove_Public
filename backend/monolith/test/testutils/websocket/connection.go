package testws

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func CreateWSConn(t testing.TB, role string, id int64, httpAddr string, dialer *websocket.Dialer) *websocket.Conn {
	t.Helper()
	connURL := fmt.Sprintf("ws://%s/ws/%s?id=%d", httpAddr, role, id)
	conn, res, err := dialer.Dial(connURL, http.Header{})
	if err != nil {
		if res != nil && res.Body != nil {
			bReason, _ := io.ReadAll(res.Body)
			t.Fatalf("dial failed: %v, reason: %v\n", err, string(bReason))
		}
		t.Fatalf("Client %d: Failed to connect: %v", id, err)
	}

	return conn
}
