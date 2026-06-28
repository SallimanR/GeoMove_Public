package command

import (
	"context"
	"fmt"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/domain/repository"
)

type UpdateDriverRealtimeCommand struct {
	DriverID    uint32
	Coordinates []dto.LocationRaw
	Time        uint64
	Distance    float32
}

type UpdateDriverRealtimeHandler struct {
	repo repository.GeolocationRepository
}

func NewUpdateDriverRealtimeHandler(repo repository.GeolocationRepository) *UpdateDriverRealtimeHandler {
	return &UpdateDriverRealtimeHandler{
		repo: repo,
	}
}

func (h *UpdateDriverRealtimeHandler) Handle(ctx context.Context, cmd UpdateDriverRealtimeCommand) error {
	// var isStaying bool
	// if len(cmd.Coordinates) < 2 || cmd.Distance < 5 {
	// 	return nil
	// }

	queryParams := entity.NewDriverRealtime(
		cmd.DriverID,
		*entity.ArrayToPoints(cmd.Coordinates),
		float32(cmd.Time),
		cmd.Distance,
	)

	// TODO: implement retry logic
	err := h.repo.UpdateDriverRealtime(ctx, queryParams)
	if err != nil {
		return fmt.Errorf("failed to update location in DB: %s", err)
	}

	return nil
}

func arrayToLocationPoints(coords []dto.LocationRaw) *[]dto.Location {
	points := make([]dto.Location, len(coords))
	for i, coord := range coords {
		points[i].Latitude = coord[0]
		points[i].Longitude = coord[1]
	}
	return &points
}
