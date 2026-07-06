package repository

import (
	"context"
	"monolith/internal/domains/driver/domain/entity"
)

type DriverFilter struct {
	UserLocation entity.Location
	WorkStarts   *string
	WorkEnds     *string
	MinRating    *float32
}

type DriverRepository interface {
	// Write (commands)
	CreateDriver(ctx context.Context, driver *entity.Driver) (entity.DriverID, error)

	// Read (queries)
	GetDriverByID(ctx context.Context, id entity.DriverID) (*entity.Driver, error)
	GetFilteredDrivers(ctx context.Context, filter DriverFilter) ([]entity.Driver, error)
}
