package query

import (
	"context"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
)

type GetOrderByIDHandler struct {
	repo repository.OrderRepository
}

func NewGetOrderByIDHandler(repo repository.OrderRepository) *GetOrderByIDHandler {
	return &GetOrderByIDHandler{repo: repo}
}

func (h *GetOrderByIDHandler) Handle(ctx context.Context, id int64) (*entity.Order, error) {
	return h.repo.GetOrderByID(ctx, id)
}
