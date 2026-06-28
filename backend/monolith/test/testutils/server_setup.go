package testutils

import (
	"testing"
	"time"

	"monolith/cmd/server"
)

func NewServer(t testing.TB, opts ...server.Option) *server.Server {
	srv, err := server.NewServer(opts...)
	if err != nil {
		t.Fatalf("Failed to setup server: %s", err)
	}

	return srv
}

func StartServer(t testing.TB, srv *server.Server) {
	err := srv.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %s", err)
	}
	time.Sleep(100 * time.Millisecond)
}
