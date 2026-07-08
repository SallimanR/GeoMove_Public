package postgres_test

import (
	"context"
	"log"
	"testing"

	"monolith/internal/domains/driver/application/commands"
	"monolith/internal/domains/driver/application/queries"
	"monolith/internal/domains/driver/domain/entities"
	"monolith/pkg/logger"
	"monolith/test/testutils"
	"monolith/test/testutils/domains"
)

func TestDriverRepository(t *testing.T) {
	ctx := context.Background()
	db, dbName, dbCleaunp := testutils.CreateTestDB(t, ctx, adminConn, adminConnString, templateDBName)
	log.Println("dbName: ", dbName)
	defer dbCleaunp()

	driverDomain := domains.NewDriverDomain(t, db)

	logger := logger.NewDebugLogger()

	var driver *entities.Driver
	t.Run("CreateDriver", func(t *testing.T) {
		createDriverCmd := commands.CreateDriverCommand{
			Latitude:  38.4,
			Longitude: 55.5,
		}
		var err error
		driver, err = driverDomain.Commands.CreateDriver.Handle(ctx, createDriverCmd)
		if err != nil {
			t.Fatalf("failed to create driver: %s", err)
		}
	})

	t.Run("GetDriverByID", func(t *testing.T) {
		findQuery := queries.GetDriverByIDQuery{DriverID: driver.ID}
		dbDriver, err := driverDomain.Queries.FindDriverByID.Handle(ctx, findQuery)
		if err != nil {
			t.Fatalf("failed to find driver by id: %s", err)
		}
		logger.Debug().Any("driver", dbDriver).Send()
	})
}
