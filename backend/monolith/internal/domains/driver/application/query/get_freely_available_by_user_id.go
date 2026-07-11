package query

import (
	"context"

	"monolith/internal/domains/driver/domain/entity"
	"monolith/internal/domains/driver/domain/repository"
)

type GetFreelyAvailableByUserIDHandler struct {
	repo repository.FreelyAvailableRepository
}

func NewGetFreelyAvailableByUserIDHandler(repo repository.FreelyAvailableRepository) *GetFreelyAvailableByUserIDHandler {
	return &GetFreelyAvailableByUserIDHandler{repo: repo}
}

func (h *GetFreelyAvailableByUserIDHandler) Handle(ctx context.Context, userID int64) (*entity.FreelyAvailable, error) {
	return h.repo.GetFreelyAvailableByUserID(ctx, userID)
}
