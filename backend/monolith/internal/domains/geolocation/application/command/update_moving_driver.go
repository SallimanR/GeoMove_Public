package command

import (
	"context"
	"fmt"
	"time"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
)

type UpdateMovingDriverCommand struct {
	DriverID    int64
	Coordinates []dto.LocationRaw
	TravelTime  time.Time
	PathMeters  uint32
}

type UpdateMovingDriverHandler struct {
	repo repository.GeolocationRepository
}

func NewUpdateMovingDriverHandler(repo repository.GeolocationRepository) *UpdateMovingDriverHandler {
	return &UpdateMovingDriverHandler{
		repo: repo,
	}
}

func (h *UpdateMovingDriverHandler) Handle(ctx context.Context, cmd UpdateMovingDriverCommand) error {
	// var isStaying bool
	// if len(cmd.Coordinates) < 2 || cmd.Distance < 5 {
	// 	return nil
	// }

	queryParams := entity.NewMovingDriver(
		cmd.DriverID,
		*entity.ArrayToPoints(cmd.Coordinates),
		cmd.TravelTime,
		cmd.PathMeters,
	)

	err := h.repo.UpdateMovingDriver(ctx, queryParams)
	if err != nil {
		return fmt.Errorf("failed to update location in DB: %s", err)
	}

	return nil
}
