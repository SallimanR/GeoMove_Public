package repository

import (
	"context"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
)

type GeolocationRepository interface {
	UpdateMovingDriver(context.Context, *entity.MovingDriver) error
	DeleteStaleMovingDrivers(ctx context.Context, cutoff string) error

	GetMovingDriverByID(ctx context.Context, id int64) (*entity.MovingDriver, error)
	GetClosestWithinRadiusMovingDrivers(ctx context.Context, location dto.Location, radius_meters uint32) ([]entity.MovingDriver, error)
}
