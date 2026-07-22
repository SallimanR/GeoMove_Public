package command

import (
	"context"
	"time"

	"monolith/internal/domains/driver/domain/repository"
)

type UpdateDriverCommand struct {
	UserID              int64
	Name                string
	Phone               *string
	WorkStarts          *time.Time
	WorkEnds            *time.Time
	Latitude            float32
	Longitude           float32
	MaxCarWeightKg      *int32
	MaxCarLengthMeters  *float32
	Address             string
}

type UpdateDriverHandler struct {
	repo repository.DriverRepository
}

func NewUpdateDriverHandler(repo repository.DriverRepository) *UpdateDriverHandler {
	return &UpdateDriverHandler{repo: repo}
}

func (h *UpdateDriverHandler) Handle(ctx context.Context, cmd UpdateDriverCommand) error {
	existing, err := h.repo.GetDriverByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	existing.Name = cmd.Name
	if cmd.Phone != nil {
		existing.Phone = *cmd.Phone
	} else {
		existing.Phone = ""
	}
	existing.WorkStarts = cmd.WorkStarts
	existing.WorkEnds = cmd.WorkEnds
	existing.Address = cmd.Address

	if err := existing.UpdateLocation(cmd.Latitude, cmd.Longitude); err != nil {
		return err
	}

	if err := h.repo.UpdateDriver(ctx, existing); err != nil {
		return err
	}

	if cmd.MaxCarWeightKg != nil && cmd.MaxCarLengthMeters != nil {
		if err := h.repo.UpsertTowDriver(ctx, cmd.UserID, *cmd.MaxCarWeightKg, *cmd.MaxCarLengthMeters); err != nil {
			return err
		}
	}

	return nil
}
