package entity

import (
	"fmt"
	"testing"

	"monolith/pkg/logger"
)

func TestCreateDriver(t *testing.T) {
	logger := logger.NewDebugLogger()

	testCases := []struct {
		driverOpts DriverOptions
		wantErr    bool
		errMsg     string
	}{
		{
			driverOpts: DriverOptions{
				Location:   Location{Lat: 55, Lon: 38},
				WorkStarts: nil,
				WorkEnds:   nil,
			},
			wantErr: false,
			errMsg:  "driver with allowed coordinates treated as invalid",
		},
		{
			driverOpts: DriverOptions{
				Location:   Location{Lat: 30, Lon: -20},
				WorkStarts: nil,
				WorkEnds:   nil,
			},
			wantErr: true,
			errMsg:  "driver with not allowed coordinates treated as valid",
		},
	}

	for i, tc := range testCases {
		driver, err := NewDriver(tc.driverOpts)

		logger.Info().AnErr("error", err).Send()
		logger.Info().Any(fmt.Sprintf("driver %d", i), driver).Send()
		if !tc.wantErr {
			if err != nil {
				t.Fatalf("%s: %s", tc.errMsg, err)
			}
		} else {
			if err == nil {
				t.Fatalf("%s", tc.errMsg)
			}
		}
	}
}
