package testutils

import (
	"io"
	"net/http"
	"testing"
)

func PingGraphhopper(t testing.TB) {
	t.Helper()

	const graphhopperURL = "http://localhost:8989/health"
	resp, err := http.Get(graphhopperURL)
	if err != nil {
		t.Fatalf("failed to ping graphhopper: %s", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if string(respBody) != "OK" {
		t.Fatalf("graphhopper is unhealthy")
	}
}
