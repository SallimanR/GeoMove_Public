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
	CreateTowDriver(ctx context.Context, driverID int64, maxWeightKg int32, maxLengthM float32) error
	UpdateProfileImage(ctx context.Context, userID int64, imageURL string) error
	UpdateDriver(ctx context.Context, driver *entity.Driver) error
	UpsertTowDriver(ctx context.Context, driverID int64, maxWeightKg int32, maxLengthM float32, carPhotoMain string, carPhotos *string) error

	GetDriverByUserID(ctx context.Context, userID int64) (*entity.Driver, error)
	GetFilteredDrivers(ctx context.Context, filter DriverFilter) ([]entity.Driver, error)
}
