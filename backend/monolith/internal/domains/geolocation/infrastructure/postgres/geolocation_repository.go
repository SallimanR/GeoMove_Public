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

func (r *GeolocationRepository) UpdateMovingDriver(ctx context.Context, gps *entity.MovingDriver) error {
	return r.queries.UpdateMovingDriver(ctx, sqlc.UpdateMovingDriverParams{
		DriverID:   gps.DriverID,
		TravelTime: gps.TravelTime,
		PathMeters: gps.PathMeters,
		Lon:        gps.Longitude,
		Lat:        gps.Latitude,
	})
}

func (r *GeolocationRepository) DeleteStaleMovingDrivers(ctx context.Context, cutoff string) error {
	return r.queries.DeleteStaleMovingDrivers(ctx, cutoff)
}

func (r *GeolocationRepository) GetMovingDriverByID(ctx context.Context, id int64) (*entity.MovingDriver, error) {
	result, err := r.queries.GetMovingDriverByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.MovingDriver{
		DriverID:   result.DriverID,
		Latitude:   result.Lat,
		Longitude:  result.Lon,
		TravelTime: result.TravelTime,
		PathMeters: result.PathMeters,
	}, nil
}

func (r *GeolocationRepository) GetClosestWithinRadiusMovingDrivers(ctx context.Context, location dto.Location, radiusMeters uint32) ([]entity.MovingDriver, error) {
	result, err := r.queries.GetClosestWithinRadiusMovingDriver(ctx, sqlc.GetClosestWithinRadiusMovingDriverParams{
		Lon:    float32(location.Longitude),
		Lat:    float32(location.Latitude),
		Radius: int32(radiusMeters),
	})
	if err != nil {
		return nil, err
	}
	output := make([]entity.MovingDriver, len(result))
	for i, row := range result {
		output[i] = entity.MovingDriver{
			DriverID:   row.DriverID,
			Latitude:   row.Lat,
			Longitude:  row.Lon,
			TravelTime: row.TravelTime,
			PathMeters: row.PathMeters,
		}
	}
	return output, nil
}
