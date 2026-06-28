package query

import (
	"context"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
)

type FindClosestDriversRealtimeQuery struct {
	Location dto.Location
}

type FindClosestDriversRealtimeHandler struct {
	repo repository.GeolocationRepository
}

func NewFindClosestDriversRealtimeHandler(repo repository.GeolocationRepository) FindClosestDriversRealtimeHandler {
	return FindClosestDriversRealtimeHandler{repo: repo}
}

func (h *FindClosestDriversRealtimeHandler) Handle(ctx context.Context, query FindClosestDriversRealtimeQuery) ([]entity.DriverDistance, error) {
	result, err := h.repo.FindClosestDriversRealtime(ctx, query.Location)
	if err != nil {
		return nil, err
	}
	return result, nil
}
