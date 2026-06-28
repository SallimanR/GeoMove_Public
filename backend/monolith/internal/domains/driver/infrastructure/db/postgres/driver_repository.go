package postgres

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
	"monolith/internal/domains/driver/infrastructure/db/sqlc"
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
	row, err := r.queries.FindDriverByID(ctx, uint32(id))
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
