package query

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type GetFilteredDriversQuery struct {
	UserLat    float32
	UserLon    float32
	WorkStarts *string
	WorkEnds   *string
	MinRating  *float32
}

type GetFilteredDriversHandler struct {
	driverRepo repository.DriverRepository
}

func NewGetFilteredDriversHandler(driverRepo repository.DriverRepository) *GetFilteredDriversHandler {
	return &GetFilteredDriversHandler{
		driverRepo: driverRepo,
	}
}

func (h *GetFilteredDriversHandler) Handle(ctx context.Context, query GetFilteredDriversQuery) ([]entity.Driver, error) {
	params := repository.DriverFilter{
		UserLocation: entity.Location{
			Lat: query.UserLat,
			Lon: query.UserLon,
		},
		WorkStarts: query.WorkStarts,
		WorkEnds:   query.WorkEnds,
		MinRating:  query.MinRating,
	}

	drivers, err := h.driverRepo.GetFilteredDrivers(ctx, params)
	if err != nil {
		return nil, err
	}
	return drivers, nil
}
