package command

import (
	"context"

	"monolith/internal/domains/driver/domain/repository"
)

type DeleteFreelyAvailableCommand struct {
	UserID int64
}

type DeleteFreelyAvailableHandler struct {
	repo repository.FreelyAvailableRepository
}

func NewDeleteFreelyAvailableHandler(repo repository.FreelyAvailableRepository) *DeleteFreelyAvailableHandler {
	return &DeleteFreelyAvailableHandler{repo: repo}
}

func (h *DeleteFreelyAvailableHandler) Handle(ctx context.Context, cmd DeleteFreelyAvailableCommand) error {
	return h.repo.DeleteFreelyAvailable(ctx, cmd.UserID)
}
