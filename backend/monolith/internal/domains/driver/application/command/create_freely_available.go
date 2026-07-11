package command

import (
	"context"
	"time"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
	"monolith/pkg/geo"
)

type CreateFreelyAvailableCommand struct {
	UserID       int64
	FromDate     time.Time
	ToDate       time.Time
	FromLocation entity.LocationWithAddress
	ToLocations  []entity.LocationWithAddress
	EnRouteOrder *bool
	TariffPerKm  *float32
}

type CreateFreelyAvailableHandler struct {
	repo repository.FreelyAvailableRepository
}

func NewCreateFreelyAvailableHandler(repo repository.FreelyAvailableRepository) *CreateFreelyAvailableHandler {
	return &CreateFreelyAvailableHandler{repo: repo}
}

func (h *CreateFreelyAvailableHandler) Handle(ctx context.Context, cmd CreateFreelyAvailableCommand) error {
	fromLoc := resolveAddress(ctx, cmd.FromLocation)
	toLocs := make([]entity.LocationWithAddress, 0, len(cmd.ToLocations))
	for _, l := range cmd.ToLocations {
		toLocs = append(toLocs, resolveAddress(ctx, l))
	}
	fa := entity.NewFreelyAvailable(cmd.UserID, cmd.FromDate, cmd.ToDate, fromLoc, toLocs, cmd.EnRouteOrder, cmd.TariffPerKm)
	return h.repo.CreateFreelyAvailable(ctx, fa)
}

func resolveAddress(ctx context.Context, loc entity.LocationWithAddress) entity.LocationWithAddress {
	addr, err := geo.ReverseGeocode(ctx, loc.Lat, loc.Lon)
	if err != nil {
		addr = ""
	}
	return entity.LocationWithAddress{Lat: loc.Lat, Lon: loc.Lon, Address: addr}
}
