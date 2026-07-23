package command

import (
	"context"
	"time"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type CreateDriverCommand struct {
	UserID             int64
	Name               string
	Phone              *string
	WorkStarts         *time.Time
	WorkEnds           *time.Time
	Latitude           float32
	Longitude          float32
	MaxCarWeightKg     *int32
	MaxCarLengthMeters *float32
	Address            string
	CarPhotoMain       string
	CarPhotos          *string
}

type CreateDriverHandler struct {
	repo repository.DriverRepository
}

func NewCreateDriverHandler(repo repository.DriverRepository) *CreateDriverHandler {
	return &CreateDriverHandler{
		repo: repo,
	}
}

func (h *CreateDriverHandler) Handle(ctx context.Context, cmd CreateDriverCommand) error {
	driverOpts := entity.DriverOptions{
		UserID:     cmd.UserID,
		Name:       cmd.Name,
		Phone:      cmd.Phone,
		WorkStarts: cmd.WorkStarts,
		WorkEnds:   cmd.WorkEnds,
		Location:   entity.Location{Lat: cmd.Latitude, Lon: cmd.Longitude},
		Address:    cmd.Address,
	}
	driver, err := entity.NewDriver(driverOpts)
	if err != nil {
		return err
	}
	err = h.repo.CreateDriver(ctx, driver)
	if err != nil {
		return err
	}

	if cmd.MaxCarWeightKg != nil && cmd.MaxCarLengthMeters != nil {
		if err := h.repo.UpsertTowDriver(ctx, driver.UserID, *cmd.MaxCarWeightKg, *cmd.MaxCarLengthMeters, cmd.CarPhotoMain, cmd.CarPhotos); err != nil {
			return err
		}
	}

	return nil
}
