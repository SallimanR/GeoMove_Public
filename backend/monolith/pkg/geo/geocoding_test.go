package geo

import (
	"context"
	"testing"

	"monolith/pkg/logger"
)

func Test_ReverseGeocodeCountry(t *testing.T) {
	logger := logger.NewDebugLogger()

	countryCode, err := ResolveCountry(context.Background(), 55, 38)
	if err != nil {
	}
	logger.Info().Str("countryCode", countryCode).Send()
	logger.Info().AnErr("error", err).Send()

	countryCode, err = ResolveCountry(context.Background(), -30, 20)
	if err != nil {
	}
	logger.Info().Str("countryCode", countryCode).Send()
	logger.Info().AnErr("error", err).Send()
}
