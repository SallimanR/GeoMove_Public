package query

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type GetDriverByUserIDQuery struct {
	UserID int64
}

type GetDriverByUserIDHandler struct {
	driverRepo repository.DriverRepository
}

func NewGetDriverByUserIDHandler(driverRepo repository.DriverRepository) *GetDriverByUserIDHandler {
	return &GetDriverByUserIDHandler{
		driverRepo: driverRepo,
	}
}

func (h *GetDriverByUserIDHandler) Handle(ctx context.Context, query GetDriverByUserIDQuery) (*entity.Driver, error) {
	driver, err := h.driverRepo.GetDriverByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	return driver, nil
}
