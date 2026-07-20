package testutils

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"monolith/internal/domains/geolocation/infrastructure/graphopper"
)

func PingGraphhopper(t testing.TB) {
	t.Helper()

	graphhopperURL := fmt.Sprintf("%s/health", graphopper.APIBase)
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
