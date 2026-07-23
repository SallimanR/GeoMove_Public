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
	phone := &driver.Phone
	if driver.Phone == "" {
		phone = nil
	}
	rating := driver.Rating
	address := &driver.Address
	if driver.Address == "" {
		address = nil
	}
	query := sqlc.CreateDriverParams{
		UserID:     driver.UserID,
		Name:       driver.Name,
		Phone:      phone,
		WorkStarts: driver.WorkStarts,
		WorkEnds:   driver.WorkEnds,
		Rating:     &rating,
		Address:    address,
		Lat:        driver.Location.Lat,
		Lon:        driver.Location.Lon,
	}
	err := r.queries.CreateDriver(ctx, query)

	return err
}

func (r *DriverRepository) UpdateProfileImage(ctx context.Context, userID int64, imageURL string) error {
	return r.queries.UpdateDriverProfileImage(ctx, sqlc.UpdateDriverProfileImageParams{
		UserID:       userID,
		ProfileImage: &imageURL,
	})
}

func (r *DriverRepository) GetDriverByUserID(ctx context.Context, userID int64) (*entity.Driver, error) {
	row, err := r.queries.GetDriverByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	phone := ""
	if row.Phone != nil {
		phone = *row.Phone
	}
	rating := float32(0)
	if row.Rating != nil {
		rating = *row.Rating
	}

	driver := &entity.Driver{
		UserID:             row.UserID,
		Name:               row.Name,
		Phone:              phone,
		ProfileImage:       row.ProfileImage,
		WorkStarts:         row.WorkStarts,
		WorkEnds:           row.WorkEnds,
		IsAvailable:        row.IsAvailable,
		LastSeen:           row.LastSeen.Time,
		Rating:             rating,
		Location:           entity.Location{Lat: row.Lat, Lon: row.Lon},
		MaxCarWeightKg:     row.MaxCarWeightKg,
		MaxCarLengthMeters: row.MaxCarLengthMeters,
		Address:            row.Address,
		CarPhotoMain:       row.CarPhotoMain,
		CarPhotos:          row.CarPhotos,
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
		phone := ""
		if row.Phone != nil {
			phone = *row.Phone
		}
		rating := float32(0)
		if row.Rating != nil {
			rating = *row.Rating
		}
		resp = append(resp, entity.Driver{
			UserID:             row.UserID,
			Name:               row.Name,
			Phone:              phone,
			ProfileImage:       row.ProfileImage,
			WorkStarts:         row.WorkStarts,
			WorkEnds:           row.WorkEnds,
			IsAvailable:        row.IsAvailable,
			LastSeen:           row.LastSeen.Time,
			Rating:             rating,
			Location:           entity.Location{Lat: row.Lat, Lon: row.Lon},
			MaxCarWeightKg:     row.MaxCarWeightKg,
			MaxCarLengthMeters: row.MaxCarLengthMeters,
			Address:            row.Address,
			CarPhotoMain:       row.CarPhotoMain,
			CarPhotos:          row.CarPhotos,
		})
	}

	return resp, nil
}

func (r *DriverRepository) CreateTowDriver(ctx context.Context, driverID int64, maxWeightKg int32, maxLengthM float32) error {
	return r.queries.CreateTowDriver(ctx, sqlc.CreateTowDriverParams{
		DriverID:           driverID,
		MaxCarWeightKg:     maxWeightKg,
		MaxCarLengthMeters: maxLengthM,
	})
}

func (r *DriverRepository) UpdateDriver(ctx context.Context, driver *entity.Driver) error {
	phone := &driver.Phone
	if driver.Phone == "" {
		phone = nil
	}
	address := &driver.Address
	if driver.Address == "" {
		address = nil
	}
	return r.queries.UpdateDriver(ctx, sqlc.UpdateDriverParams{
		UserID:     driver.UserID,
		Name:       driver.Name,
		Phone:      phone,
		WorkStarts: driver.WorkStarts,
		WorkEnds:   driver.WorkEnds,
		Column6:    driver.Location.Lon,
		Column7:    driver.Location.Lat,
		Address:    address,
	})
}

func (r *DriverRepository) UpsertTowDriver(ctx context.Context, driverID int64, maxWeightKg int32, maxLengthM float32, carPhotoMain string, carPhotos *string) error {
	return r.queries.UpsertTowDriver(ctx, sqlc.UpsertTowDriverParams{
		DriverID:           driverID,
		MaxCarWeightKg:     maxWeightKg,
		MaxCarLengthMeters: maxLengthM,
		CarPhotoMain:       carPhotoMain,
		CarPhotos:          carPhotos,
	})
}
