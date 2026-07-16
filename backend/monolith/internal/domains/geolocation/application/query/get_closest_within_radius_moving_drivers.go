package query

import (
	"context"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
)

type GetClosestWithinRadiusMovingDriversQuery struct {
	Location     dto.Location
	RadiusMeters uint32
}

type GetClosestWithinRadiusMovingDriversHandler struct {
	repo repository.GeolocationRepository
}

func NewGetClosestWithinRadiusMovingDriversHandler(repo repository.GeolocationRepository) GetClosestWithinRadiusMovingDriversHandler {
	return GetClosestWithinRadiusMovingDriversHandler{repo: repo}
}

func (h *GetClosestWithinRadiusMovingDriversHandler) Handle(ctx context.Context, query GetClosestWithinRadiusMovingDriversQuery) ([]entity.MovingDriver, error) {
	result, err := h.repo.GetClosestWithinRadiusMovingDrivers(ctx, query.Location, query.RadiusMeters)
	if err != nil {
		return nil, err
	}
	return result, nil
}
