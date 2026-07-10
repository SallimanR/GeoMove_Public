package command

import (
	"context"

	"monolith/internal/domains/driver/domain/repository"
)

type UpdateProfileImageCommand struct {
	UserID   int64
	ImageURL string
}

type UpdateProfileImageHandler struct {
	repo repository.DriverRepository
}

func NewUpdateProfileImageHandler(repo repository.DriverRepository) *UpdateProfileImageHandler {
	return &UpdateProfileImageHandler{repo: repo}
}

func (h *UpdateProfileImageHandler) Handle(ctx context.Context, cmd UpdateProfileImageCommand) error {
	return h.repo.UpdateProfileImage(ctx, cmd.UserID, cmd.ImageURL)
}
