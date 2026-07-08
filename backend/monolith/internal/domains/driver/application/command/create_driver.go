package command

import (
	"context"
	"time"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type CreateDriverCommand struct {
	UserID     int64
	Name       string
	WorkStarts *time.Time
	WorkEnds   *time.Time
	Latitude   float32
	Longitude  float32
}

type CreateDriverHandler struct {
	repo repository.DriverRepository
}

func NewCreateDriverHandler(repo repository.DriverRepository) *CreateDriverHandler {
	return &CreateDriverHandler{
		repo: repo,
	}
}

func (h *CreateDriverHandler) Handle(ctx context.Context, cmd CreateDriverCommand) (*entity.Driver, error) {
	driverOpts := entity.DriverOptions{
		UserID:     cmd.UserID,
		Name:       cmd.Name,
		WorkStarts: cmd.WorkStarts,
		WorkEnds:   cmd.WorkEnds,
		Location:   entity.Location{Lat: cmd.Latitude, Lon: cmd.Longitude},
	}
	driver, err := entity.NewDriver(driverOpts)
	if err != nil {
		return nil, err
	}
	driverID, err := h.repo.CreateDriver(ctx, driver)
	driver.ID = driverID
	if err != nil {
		return nil, err
	}
	return driver, nil
}
