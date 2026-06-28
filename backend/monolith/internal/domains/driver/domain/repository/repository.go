package repository

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
)

type DriverRepository interface {
	// Write (commands)
	CreateDriver(ctx context.Context, driver *entity.Driver) (entity.DriverID, error)

	// Read (queries)
	GetDriverByID(ctx context.Context, id entity.DriverID) (*entity.Driver, error)
}
