package graphopper

import (
	"testing"
	"unsafe"

	"monolith/internal/domains/geolocation/interface/websocket/pb"
	"monolith/pkg/logger"
)

func TestMapMatch(t *testing.T) {
	logger := logger.NewDebugLogger()

	// TODO: both: in-city; inter-city
	testCases := []struct {
		name    string
		data    pb.LocationUpdate
		wantErr bool
	}{
		{
			name: "test moscow",
			data: pb.LocationUpdate{
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
			data: pb.LocationUpdate{
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

			// TODO: measure matched?
			matched, ok := MatchMapCoordinates(&tc.data)
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
