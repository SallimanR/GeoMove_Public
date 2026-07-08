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
	CreateDriver(ctx context.Context, driver *entity.Driver) error

	GetDriverByUserID(ctx context.Context, userID int64) (*entity.Driver, error)
	GetFilteredDrivers(ctx context.Context, filter DriverFilter) ([]entity.Driver, error)
}
