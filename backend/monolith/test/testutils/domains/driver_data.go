package domains

import (
	"context"
	"testing"
	"time"

	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/domain/entity"
)

// func CreateDriver(t testing.TB, domain *DriverDomain, cmd commands.CreateDriverCommand) *entities.Driver {
func CreateDriver(t testing.TB, handler *command.CreateDriverHandler, cmd command.CreateDriverCommand) *entity.Driver {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// driver, err := domain.Commands.CreateDriver.Handle(ctx, cmd)
	driver, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("failed to create driver: %s", err)
	}
	return driver
}
