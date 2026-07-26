package command

import (
	"context"
	"fmt"

	"monolith/internal/domains/order/domain/repository"
)

type DeclineOrderHandler struct {
	repo repository.OrderRepository
}

func NewDeclineOrderHandler(repo repository.OrderRepository) *DeclineOrderHandler {
	return &DeclineOrderHandler{repo: repo}
}

func (h *DeclineOrderHandler) Handle(ctx context.Context, orderID, driverID int64) error {
	if err := h.repo.DeclineOrder(ctx, orderID, driverID); err != nil {
		return fmt.Errorf("ошибка отказа от заказа: %w", err)
	}
	return nil
}
