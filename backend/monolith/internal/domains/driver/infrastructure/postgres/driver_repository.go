package postgres

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
	"monolith/internal/domains/driver/infrastructure/postgres/sqlc"
	"monolith/pkg/db"
)

type DriverRepository struct {
	queries sqlc.Queries
}

func NewDriverRepository(queries *sqlc.Queries) repository.DriverRepository {
	return &DriverRepository{
		queries: *queries,
	}
}

func (r *DriverRepository) CreateDriver(ctx context.Context, driver *entity.Driver) error {
	query := sqlc.CreateDriverParams{
		UserID:     driver.UserID,
		Name:       driver.Name,
		WorkStarts: driver.WorkStarts,
		WorkEnds:   driver.WorkStarts,
		Lat:        driver.Location.Lat,
		Lon:        driver.Location.Lon,
	}
	err := r.queries.CreateDriver(ctx, query)

	return err
}

func (r *DriverRepository) GetDriverByUserID(ctx context.Context, userID int64) (*entity.Driver, error) {
	row, err := r.queries.GetDriverByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	driver := &entity.Driver{
		UserID:       row.UserID,
		Name:         row.Name,
		ProfileImage: row.ProfileImage,
		WorkStarts:   row.WorkStarts,
		WorkEnds:     row.WorkEnds,
		IsAvailable:  row.IsAvailable,
		LastSeen:     row.LastSeen.Time,
		Rating:       row.Rating,
		Location:     entity.Location{Lat: row.Lat, Lon: row.Lon},
	}

	return driver, nil
}

func (r *DriverRepository) GetFilteredDrivers(ctx context.Context, filter repository.DriverFilter) ([]entity.Driver, error) {
	params := sqlc.GetFilteredDriversParams{
		Lat:       filter.UserLocation.Lat,
		Lon:       filter.UserLocation.Lon,
		MinRating: filter.MinRating,
	}

	if filter.WorkStarts != nil {
		params.WorkStarts = db.StringToPgTime(*filter.WorkStarts)
	}
	if filter.WorkEnds != nil {
		params.WorkEnds = db.StringToPgTime(*filter.WorkEnds)
	}

	rows, err := r.queries.GetFilteredDrivers(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := make([]entity.Driver, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, entity.Driver{
			UserID:       row.UserID,
			Name:         row.Name,
			ProfileImage: row.ProfileImage,
			WorkStarts:   row.WorkStarts,
			WorkEnds:     row.WorkEnds,
			IsAvailable:  row.IsAvailable,
			LastSeen:     row.LastSeen.Time,
			Rating:       row.Rating,
			Location:     entity.Location{Lat: row.Lat, Lon: row.Lon},
		})
	}

	return resp, nil
}
