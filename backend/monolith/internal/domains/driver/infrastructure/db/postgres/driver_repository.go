package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
	"monolith/internal/domains/driver/infrastructure/db/sqlc"
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

func (r *DriverRepository) CreateDriver(ctx context.Context, driver *entity.Driver) (entity.DriverID, error) {
	query := sqlc.CreateDriverParams{
		// WorkStarts: driver.WorkStarts,
		Column3: driver.Location.Lat,
		Column4: driver.Location.Lon,
	}
	driverID, err := r.queries.CreateDriver(ctx, query)
	if err != nil {
		return 0, err
	}
	return entity.DriverID(driverID), nil
}

func (r *DriverRepository) GetDriverByID(ctx context.Context, id entity.DriverID) (*entity.Driver, error) {
	row, err := r.queries.GetDriverByID(ctx, uint32(id))
	if err != nil {
		return nil, err
	}

	driver := &entity.Driver{
		ID: entity.DriverID(row.ID),
		// WorkStarts: row.WorkStarts,
		Location: entity.Location{Lat: row.Lat, Lon: row.Lon},
	}

	return driver, nil
}

func (r *DriverRepository) GetFilteredDrivers(ctx context.Context, filter repository.DriverFilter) ([]entity.Driver, error) {
	params := sqlc.GetFilteredDriversParams{
		Lat: filter.UserLocation.Lat,
		Lon: filter.UserLocation.Lon,
	}

	if filter.WorkStarts != nil {
		params.WorkStarts = db.StringToPgTime(*filter.WorkStarts)
	} else {
		params.WorkStarts = pgtype.Time{Valid: false}
	}
	if filter.WorkEnds != nil {
		params.WorkEnds = db.StringToPgTime(*filter.WorkEnds)
	} else {
		params.WorkEnds = pgtype.Time{Valid: false}
	}
	if filter.MinRating != nil {
		params.MinRating = pgtype.Float4{Float32: *filter.MinRating, Valid: true}
	} else {
		params.MinRating = pgtype.Float4{Valid: false}
	}

	rows, err := r.queries.GetFilteredDrivers(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := make([]entity.Driver, 0, len(rows))
	for _, row := range rows {
		driverID := entity.DriverID(row.ID)

		var workStarts *time.Time
		if row.WorkStarts.Valid {
			t := time.UnixMicro(row.WorkStarts.Microseconds)
			workStarts = &t
		}

		var workEnds *time.Time
		if row.WorkEnds.Valid {
			t := time.UnixMicro(row.WorkEnds.Microseconds)
			workEnds = &t
		}

		rating := row.Rating.Float32

		resp = append(resp, entity.Driver{
			ID:         driverID,
			Name:       row.Name,
			WorkStarts: workStarts,
			WorkEnds:   workEnds,
			Rating:     &rating,
			Location: entity.Location{
				Lat: row.Lat,
				Lon: row.Lon,
			},
		})
	}

	return resp, nil
}
