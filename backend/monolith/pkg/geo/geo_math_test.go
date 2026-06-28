package geo

import (
	"strconv"
	"testing"

	"monolith/pkg/logger"
)

func TestDistance(t *testing.T) {
	logger := logger.NewDebugLogger()

	testCases := []struct {
		lat1 float64
		lon1 float64
		lat2 float64
		lon2 float64
	}{
		{
			lat1: 55.821624, lon1: 37.6354189,
			lat2: 55.742349, lon2: 37.617037,
		},
		{
			lat1: 55.751624, lon1: 37.6054189,
			lat2: 55.742349, lon2: 37.617037,
		},
		{
			lat1: 55.713652, lon1: 37.6213846,
			lat2: 55.742349, lon2: 37.617037,
		},
	}

	for i, tc := range testCases {
		distance := CalculateDistance(tc.lat1, tc.lon1, tc.lat2, tc.lon2)
		logger.Debug().Float32(strconv.Itoa(i), distance).Send()
	}
}
