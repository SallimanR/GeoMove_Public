package network

import (
	"net"
	"testing"
)

func GetFreePort(tb testing.TB) int {
	tb.Helper()
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		tb.Fatalf("failed to resolve TCP addr: %s", err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		tb.Fatalf("failed to listen to TCP: %s", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return port
}
