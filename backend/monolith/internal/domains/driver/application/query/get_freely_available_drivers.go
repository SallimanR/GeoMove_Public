package query

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type GetFreelyAvailableDriversQuery struct {
	UserLat      float32
	UserLon      float32
	EnRouteOrder *bool
	MinTariff    *float32
	MaxTariff    *float32
}

type GetFreelyAvailableDriversHandler struct {
	repo repository.FreelyAvailableRepository
}

func NewGetFreelyAvailableDriversHandler(repo repository.FreelyAvailableRepository) *GetFreelyAvailableDriversHandler {
	return &GetFreelyAvailableDriversHandler{repo: repo}
}

func (h *GetFreelyAvailableDriversHandler) Handle(ctx context.Context, qry GetFreelyAvailableDriversQuery) ([]entity.FreelyAvailableDriver, error) {
	filter := repository.FreelyAvailableFilter{
		UserLat:      qry.UserLat,
		UserLon:      qry.UserLon,
		EnRouteOrder: qry.EnRouteOrder,
		MinTariff:    qry.MinTariff,
		MaxTariff:    qry.MaxTariff,
	}
	return h.repo.GetFreelyAvailableDrivers(ctx, filter)
}
