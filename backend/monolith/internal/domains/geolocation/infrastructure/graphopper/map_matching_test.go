package graphopper

import (
	"context"
	"testing"
	"time"
	"unsafe"

	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/pkg/logger"
)

func TestMapMatch(t *testing.T) {
	logger := logger.NewDebugLogger()

	testCases := []struct {
		name    string
		data    pb.GPSUpdate
		wantErr bool
	}{
		{
			name: "test moscow",
			data: pb.GPSUpdate{
				Coordinates: []*pb.Point{
					{Latitude: 55.746055, Longitude: 37.621070},
					{Latitude: 55.745930, Longitude: 37.620733},
					{Latitude: 55.746141, Longitude: 37.620314},
					{Latitude: 55.746582, Longitude: 37.619739},
				},
				Timestamps: []int64{10, 15, 20},
			},
		},
		{
			name: "random coords with broken sequence",
			data: pb.GPSUpdate{
				Coordinates: []*pb.Point{
					{Latitude: 55.0, Longitude: 38.3},
					{Latitude: 55.2, Longitude: 38.5},
					{Latitude: 55.6, Longitude: 38.6},
				},
				Timestamps: []int64{10, 15, 20},
			},
			wantErr: true,
		},
	}
	for i := 0; i < len(testCases) && !t.Failed(); i++ {
		tc := &testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			logger.Info().Str("TEST", tc.name).Send()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			// TODO: check matched
			matched, ok := MatchMapCoordinates(ctx, &tc.data)
			if !tc.wantErr {
				if !ok {
					t.Fatalf("Failed to handle location update")
				}
			} else {
				if ok {
					t.Fatalf("False data is processed")
				}
			}

			logger.Debug().Uint16("matchedSize", uint16(unsafe.Sizeof(matched))).Send()
			logger.Debug().Any("matched", matched).Send()
		})
	}
}
