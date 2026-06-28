package query

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type GetDriverByIDQuery struct {
	DriverID entity.DriverID
}

type GetDriverByIDHandler struct {
	driverRepo repository.DriverRepository
}

func NewGetDriverByIDHandler(driverRepo repository.DriverRepository) *GetDriverByIDHandler {
	return &GetDriverByIDHandler{
		driverRepo: driverRepo,
	}
}

func (h *GetDriverByIDHandler) Handle(ctx context.Context, query GetDriverByIDQuery) (*entity.Driver, error) {
	driver, err := h.driverRepo.GetDriverByID(ctx, query.DriverID)
	if err != nil {
		return nil, err
	}
	return driver, nil
}
