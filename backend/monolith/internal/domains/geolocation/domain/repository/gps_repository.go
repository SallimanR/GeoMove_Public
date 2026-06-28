package repository

import (
	"context"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
)

type GeolocationRepository interface {
	/* Write (commands) */
	CreateDriverRealtime(ctx context.Context, id uint32) error
	UpdateDriverRealtime(context.Context, *entity.DriverRealtime) error

	/* Read (queries) */
	GetDriverRealtimeByID(ctx context.Context, id uint32) (*entity.DriverRealtime, error)
	// TODO: rename to "Find" to "Get" ?
	FindClosestDriversRealtime(ctx context.Context, location dto.Location) ([]entity.DriverDistance, error)
	FindClosestDriversRealtimeH3(ctx context.Context, location dto.Location) ([]entity.DriverDistance, error)
	FindClosestWithinRadiusDriversRealtime(ctx context.Context, location dto.Location, radius uint32) ([]entity.DriverDistance, error)
}
