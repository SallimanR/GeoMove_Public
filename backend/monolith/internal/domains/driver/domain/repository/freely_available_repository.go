package repository

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
)

type FreelyAvailableFilter struct {
	UserLat     float32
	UserLon     float32
	EnRouteOrder *bool
	MinTariff   *float32
	MaxTariff   *float32
}

type FreelyAvailableRepository interface {
	CreateFreelyAvailable(ctx context.Context, fa *entity.FreelyAvailable) error
	GetFreelyAvailableByUserID(ctx context.Context, userID int64) (*entity.FreelyAvailable, error)
	UpdateFreelyAvailable(ctx context.Context, fa *entity.FreelyAvailable) error
	DeleteFreelyAvailable(ctx context.Context, userID int64) error
	GetFreelyAvailableDrivers(ctx context.Context, filter FreelyAvailableFilter) ([]entity.FreelyAvailableDriver, error)
}
