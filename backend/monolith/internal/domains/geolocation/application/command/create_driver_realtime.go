package command

import (
	"context"

	"monolith/internal/domains/geolocation/domain/repository"
)

type CreateDriverRealtimeCommand struct {
	DriverID uint32
}

type CreateDriverRealtimeHandler struct {
	repo repository.GeolocationRepository
}

func NewCreateDriverRealtimeHandler(repo repository.GeolocationRepository) *CreateDriverRealtimeHandler {
	return &CreateDriverRealtimeHandler{
		repo: repo,
	}
}

func (h *CreateDriverRealtimeHandler) Handle(ctx context.Context, cmd CreateDriverRealtimeCommand) error {
	err := h.repo.CreateDriverRealtime(ctx, cmd.DriverID)
	if err != nil {
		return err
	}
	return nil
}
