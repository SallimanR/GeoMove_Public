package query

import (
	"context"

	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
)

type GetDriverRealtimeByIDQuery struct {
	DriverID uint32
}

type GetDriverRealtimeByIDHandler struct {
	repo repository.GeolocationRepository
}

func NewGetDriverRealtimeByIDHandler(repo repository.GeolocationRepository) GetDriverRealtimeByIDHandler {
	return GetDriverRealtimeByIDHandler{repo: repo}
}

func (h *GetDriverRealtimeByIDHandler) Handle(ctx context.Context, cmd GetDriverRealtimeByIDQuery) (*entity.DriverRealtime, error) {
	result, err := h.repo.GetDriverRealtimeByID(ctx, cmd.DriverID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
