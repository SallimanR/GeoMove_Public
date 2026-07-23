package postgres

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
	"monolith/internal/domains/driver/infrastructure/postgres/sqlc"
)

type FreelyAvailableRepository struct {
	queries sqlc.Queries
}

func NewFreelyAvailableRepository(queries *sqlc.Queries) repository.FreelyAvailableRepository {
	return &FreelyAvailableRepository{
		queries: *queries,
	}
}

func (r *FreelyAvailableRepository) CreateFreelyAvailable(ctx context.Context, fa *entity.FreelyAvailable) error {
	err := r.queries.CreateFreelyAvailable(ctx, sqlc.CreateFreelyAvailableParams{
		UserID:       fa.UserID,
		FromDate:     fa.FromDate,
		ToDate:       fa.ToDate,
		FromLon:      fa.FromLocation.Lon,
		FromLat:      fa.FromLocation.Lat,
		EnRouteOrder: fa.EnRouteOrder,
		TariffPerKm:  fa.TariffPerKm,
		FromAddress:  fa.FromLocation.Address,
	})
	if err != nil {
		return err
	}

	for _, loc := range fa.ToLocations {
		err := r.queries.CreateFreelyAvailableLocation(ctx, sqlc.CreateFreelyAvailableLocationParams{
			TowDriver: fa.UserID,
			Lon:       loc.Lon,
			Lat:       loc.Lat,
			Address:   loc.Address,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FreelyAvailableRepository) GetFreelyAvailableByUserID(ctx context.Context, userID int64) (*entity.FreelyAvailable, error) {
	row, err := r.queries.GetFreelyAvailableByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	locRows, err := r.queries.GetFreelyAvailableLocations(ctx, userID)
	if err != nil {
		return nil, err
	}

	locations := make([]entity.LocationWithAddress, 0, len(locRows))
	for _, l := range locRows {
		locations = append(locations, entity.LocationWithAddress{Lat: l.Lat, Lon: l.Lon, Address: l.Address})
	}

	fa := entity.NewFreelyAvailable(
		row.UserID,
		row.FromDate,
		row.ToDate,
		entity.LocationWithAddress{Lat: row.FromLat, Lon: row.FromLon, Address: row.FromAddress},
		locations,
		row.EnRouteOrder,
		row.TariffPerKm,
	)

	return fa, nil
}

func (r *FreelyAvailableRepository) UpdateFreelyAvailable(ctx context.Context, fa *entity.FreelyAvailable) error {
	err := r.queries.UpdateFreelyAvailable(ctx, sqlc.UpdateFreelyAvailableParams{
		UserID:       fa.UserID,
		FromDate:     fa.FromDate,
		ToDate:       fa.ToDate,
		FromLon:      fa.FromLocation.Lon,
		FromLat:      fa.FromLocation.Lat,
		EnRouteOrder: fa.EnRouteOrder,
		TariffPerKm:  fa.TariffPerKm,
		FromAddress:  fa.FromLocation.Address,
	})
	if err != nil {
		return err
	}

	err = r.queries.DeleteFreelyAvailableLocations(ctx, fa.UserID)
	if err != nil {
		return err
	}

	for _, loc := range fa.ToLocations {
		err := r.queries.CreateFreelyAvailableLocation(ctx, sqlc.CreateFreelyAvailableLocationParams{
			TowDriver: fa.UserID,
			Lon:       loc.Lon,
			Lat:       loc.Lat,
			Address:   loc.Address,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FreelyAvailableRepository) DeleteFreelyAvailable(ctx context.Context, userID int64) error {
	return r.queries.DeleteFreelyAvailable(ctx, userID)
}

func (r *FreelyAvailableRepository) GetFreelyAvailableDrivers(ctx context.Context, filter repository.FreelyAvailableFilter) ([]entity.FreelyAvailableDriver, error) {
	params := sqlc.GetFreelyAvailableDriversParams{
		UserLon:      filter.UserLon,
		UserLat:      filter.UserLat,
		EnRouteOrder: filter.EnRouteOrder,
		MinTariff:    filter.MinTariff,
		MaxTariff:    filter.MaxTariff,
	}

	rows, err := r.queries.GetFreelyAvailableDrivers(ctx, params)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}

	locMap := make(map[int64][]entity.LocationWithAddress)
	if len(userIDs) > 0 {
		locRows, err := r.queries.GetFreelyAvailableLocationsByDrivers(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		for _, l := range locRows {
			locMap[l.TowDriver] = append(locMap[l.TowDriver], entity.LocationWithAddress{Lat: l.Lat, Lon: l.Lon, Address: l.Address})
		}
	}

	result := make([]entity.FreelyAvailableDriver, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.FreelyAvailableDriver{
			UserID:       row.UserID,
			FromDate:     row.FromDate,
			ToDate:       row.ToDate,
			FromLocation: entity.LocationWithAddress{Lat: row.FromLat, Lon: row.FromLon, Address: row.FromAddress},
			ToLocations:  locMap[row.UserID],
			EnRouteOrder: row.EnRouteOrder,
			TariffPerKm:  row.TariffPerKm,
			Name:         row.Name,
			Rating:       row.Rating,
			ProfileImage: row.ProfileImage,
			Distance:     row.Distance,
		})
	}

	return result, nil
}
