package domains

import (
	"context"
	"testing"
	"time"

	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/setup"
)

func CreateDriverRealtime(t testing.TB, geoDomain *setup.GeolocationDomain, cmd command.CreateDriverRealtimeCommand) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := geoDomain.Commands.CreateDriverRealtime.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("failed to create driver in realtime table: %s", err)
	}
}
