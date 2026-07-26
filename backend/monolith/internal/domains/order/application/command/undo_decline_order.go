package command

import (
	"context"
	"fmt"

	"monolith/internal/domains/order/domain/repository"
)

type UndoDeclineOrderHandler struct {
	repo repository.OrderRepository
}

func NewUndoDeclineOrderHandler(repo repository.OrderRepository) *UndoDeclineOrderHandler {
	return &UndoDeclineOrderHandler{repo: repo}
}

func (h *UndoDeclineOrderHandler) Handle(ctx context.Context, orderID, driverID int64) error {
	if err := h.repo.UndoDeclineOrder(ctx, orderID, driverID); err != nil {
		return fmt.Errorf("ошибка восстановления заказа: %w", err)
	}
	return nil
}
