package postgres

import (
	"context"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
	"monolith/internal/domains/geolocation/infrastructure/postgres/sqlc"
)

type GeolocationRepository struct {
	queries *sqlc.Queries
}

func NewGeolocationRepository(queries *sqlc.Queries) repository.GeolocationRepository {
	return &GeolocationRepository{queries: queries}
}

/* Write (commands) */

func (r *GeolocationRepository) CreateDriverRealtime(ctx context.Context, id uint32) error {
	err := r.queries.CreateDriverRealtime(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *GeolocationRepository) UpdateDriverRealtime(ctx context.Context, gps *entity.DriverRealtime) error {
	queryParams := sqlc.UpdateDriverRealtimeParams{
		DriverID: gps.DriverID,
		Column2:  gps.Latitude,
		Column3:  gps.Longitude,
		Column4:  gps.Bearing,
		Column5:  gps.AverageSpeed,
	}
	err := r.queries.UpdateDriverRealtime(ctx, queryParams)
	if err != nil {
		return err
	}
	return nil
}

/* Read (queries) */

func (r *GeolocationRepository) GetDriverRealtimeByID(ctx context.Context, id uint32) (*entity.DriverRealtime, error) {
	result, err := r.queries.GetDriverRealtimeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	domainRelatime := entity.DriverRealtime{
		DriverID:  result.DriverID,
		Latitude:  result.Lat,
		Longitude: result.Lon,
	}
	return &domainRelatime, nil
}

func (r *GeolocationRepository) FindClosestDriversRealtime(ctx context.Context, location dto.Location) ([]entity.DriverDistance, error) {
	queryParams := sqlc.FindClosestDriversRealtimeParams{
		Column1: location.Latitude,
		Column2: location.Longitude,
	}
	result, err := r.queries.FindClosestDriversRealtime(ctx, queryParams)
	if err != nil {
		return nil, err
	}
	resultsCount := len(result)
	output := make([]entity.DriverDistance, resultsCount)

	for i := 0; i < len(result); i++ {
		output[i] = entity.DriverDistance{
			DriverID:       entity.DriverID(result[i].DriverID),
			DistanceMeters: uint32(result[i].DistanceMeters),
		}
	}
	return output, err
}

func (r *GeolocationRepository) FindClosestDriversRealtimeH3(ctx context.Context, location dto.Location) ([]entity.DriverDistance, error) {
	queryParams := sqlc.FindClosestDriversRealtimeH3Params{
		Column1: location.Latitude,
		Column2: location.Longitude,
	}
	result, err := r.queries.FindClosestDriversRealtimeH3(ctx, queryParams)
	if err != nil {
		return nil, err
	}
	resultsCount := len(result)
	output := make([]entity.DriverDistance, resultsCount)

	for i := 0; i < len(result); i++ {
		output[i] = entity.DriverDistance{
			DriverID:       entity.DriverID(result[i].DriverID),
			DistanceMeters: uint32(result[i].DistanceMeters),
		}
	}

	return output, nil
}

func (r *GeolocationRepository) FindClosestWithinRadiusDriversRealtime(ctx context.Context, location dto.Location, radius uint32) ([]entity.DriverDistance, error) {
	queryParams := sqlc.FindClosestWithinRadiusDriversRealtimeParams{
		Column1: location.Latitude,
		Column2: location.Longitude,
		Column3: int32(radius),
	}
	result, err := r.queries.FindClosestWithinRadiusDriversRealtime(ctx, queryParams)
	if err != nil {
		return nil, err
	}
	resultsCount := len(result)
	output := make([]entity.DriverDistance, resultsCount)

	for i := 0; i < len(result); i++ {
		output[i] = entity.DriverDistance{
			DriverID:       entity.DriverID(result[i].DriverID),
			DistanceMeters: uint32(result[i].DistanceMeters),
		}
	}
	return output, err
}
