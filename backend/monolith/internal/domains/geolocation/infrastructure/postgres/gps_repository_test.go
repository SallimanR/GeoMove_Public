package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/setup"
	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/infrastructure/postgres"
	"monolith/internal/domains/geolocation/infrastructure/postgres/sqlc"
	"monolith/pkg/logger"
	"monolith/test/testutils"
)

func Test_FindClosest(t *testing.T) {
	logger := logger.NewLogger(zerolog.DebugLevel)

	ctx := context.Background()
	db, dbName, dbCleaunp := testutils.CreateTestDB(t, ctx, adminConn, adminConnString, templateDBName)
	defer dbCleaunp()
	logger.Info().Str("dbName: ", dbName).Send()

	geoRepo := postgres.NewGeolocationRepository(sqlc.New(db))

	driverDomain := setup.NewDriverDomain(db)

	type driverData struct {
		location dto.Location
	}
	testData := []driverData{
		{location: dto.Location{Latitude: 55.821624, Longitude: 37.6354189}},
		{location: dto.Location{Latitude: 55.751624, Longitude: 37.6054189}},
		{location: dto.Location{Latitude: 55.713652, Longitude: 37.6213846}},
	}
	for _, td := range testData {
		createDriverCmd := command.CreateDriverCommand{
			Latitude:  td.location.Latitude,
			Longitude: td.location.Longitude,
		}

		driver, err := driverDomain.Commands.CreateDriver.Handle(ctx, createDriverCmd)
		if err != nil {
			t.Fatalf("failed to create driver: %s", err)
		}

		err = geoRepo.CreateDriverRealtime(ctx, uint32(driver.ID))
		if err != nil {
			t.Fatalf("failed to create driver in realtime table: %s", err)
		}
		err = geoRepo.UpdateDriverRealtime(ctx, &entity.DriverRealtime{
			DriverID:  uint32(driver.ID),
			Latitude:  td.location.Latitude,
			Longitude: td.location.Longitude,
		})
		if err != nil {
			t.Fatalf("failed to create driver in realtime table: %s", err)
		}
	}

	testCases := []struct {
		location dto.Location
		radius   uint32
	}{
		{
			location: dto.Location{
				Latitude:  55.742349,
				Longitude: 37.617037,
			},
			radius: 8091,
		},
		{
			location: dto.Location{
				Latitude:  55.742349,
				Longitude: 37.617037,
			},
			radius: 8092,
		},
	}
	for i := 0; i < len(testCases) && !t.Failed(); i++ {
		tc := &testCases[i]

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		drivers, err := geoRepo.FindClosestDriversRealtime(ctx, tc.location)
		if err != nil {
			t.Fatalf("error finding closest drivers: %s", err)
		}
		logger.Debug().Any("FindClosestDriversRealtime", drivers).Send()

		drivers, err = geoRepo.FindClosestDriversRealtimeH3(ctx, tc.location)
		if err != nil {
			t.Fatalf("error finding closest drivers: %s", err)
		}
		logger.Debug().Any("FindClosestDriversRealtimeH3", drivers).Send()

		drivers, err = geoRepo.FindClosestWithinRadiusDriversRealtime(ctx, tc.location, tc.radius)
		if err != nil {
			t.Fatalf("error finding closest drivers: %s", err)
		}
		logger.Debug().Any("FindClosestWithinRadiusDriversRealtime", drivers).Send()
	}
}
