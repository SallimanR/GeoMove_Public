package query

import (
	"context"

	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
)

type GetMovingDriverByIDQuery struct {
	DriverID int64
}

type GetMovingDriverByIDHandler struct {
	repo repository.GeolocationRepository
}

func NewGetMovingDriverByIDHandler(repo repository.GeolocationRepository) GetMovingDriverByIDHandler {
	return GetMovingDriverByIDHandler{repo: repo}
}

func (h *GetMovingDriverByIDHandler) Handle(ctx context.Context, cmd GetMovingDriverByIDQuery) (*entity.MovingDriver, error) {
	result, err := h.repo.GetMovingDriverByID(ctx, cmd.DriverID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
