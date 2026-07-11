package command

import (
	"context"
	"time"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type UpdateFreelyAvailableCommand struct {
	UserID       int64
	FromDate     time.Time
	ToDate       time.Time
	FromLocation entity.LocationWithAddress
	ToLocations  []entity.LocationWithAddress
	EnRouteOrder *bool
	TariffPerKm  *float32
}

type UpdateFreelyAvailableHandler struct {
	repo repository.FreelyAvailableRepository
}

func NewUpdateFreelyAvailableHandler(repo repository.FreelyAvailableRepository) *UpdateFreelyAvailableHandler {
	return &UpdateFreelyAvailableHandler{repo: repo}
}

func (h *UpdateFreelyAvailableHandler) Handle(ctx context.Context, cmd UpdateFreelyAvailableCommand) error {
	fromLoc := resolveAddress(ctx, cmd.FromLocation)
	toLocs := make([]entity.LocationWithAddress, 0, len(cmd.ToLocations))
	for _, l := range cmd.ToLocations {
		toLocs = append(toLocs, resolveAddress(ctx, l))
	}
	fa := entity.NewFreelyAvailable(cmd.UserID, cmd.FromDate, cmd.ToDate, fromLoc, toLocs, cmd.EnRouteOrder, cmd.TariffPerKm)
	return h.repo.UpdateFreelyAvailable(ctx, fa)
}
