package command

import (
	"context"

	"monolith/internal/domains/geolocation/domain/repository"
)

type DeleteStaleMovingDriversHandler struct {
	repo repository.GeolocationRepository
}

func NewDeleteStaleMovingDriversHandler(repo repository.GeolocationRepository) *DeleteStaleMovingDriversHandler {
	return &DeleteStaleMovingDriversHandler{repo: repo}
}

func (h *DeleteStaleMovingDriversHandler) Handle(ctx context.Context, cutoff string) error {
	return h.repo.DeleteStaleMovingDrivers(ctx, cutoff)
}
