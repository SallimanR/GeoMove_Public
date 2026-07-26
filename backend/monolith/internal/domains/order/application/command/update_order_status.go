package command

import (
	"context"
	"fmt"
	"log"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
	"monolith/internal/notification"
)

type UpdateOrderStatusCommand struct {
	OrderID            int64
	Status             entity.OrderStatus
	DriverID           *int64
	CancellationReason *string
}

type UpdateOrderStatusHandler struct {
	repo     repository.OrderRepository
	notifSvc *notification.Service
}

func NewUpdateOrderStatusHandler(repo repository.OrderRepository, notifSvc *notification.Service) *UpdateOrderStatusHandler {
	return &UpdateOrderStatusHandler{repo: repo, notifSvc: notifSvc}
}

func (h *UpdateOrderStatusHandler) Handle(ctx context.Context, cmd UpdateOrderStatusCommand) (*entity.Order, error) {
	order, err := h.repo.GetOrderByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}

	previousStatus := order.Status

	if cmd.Status == entity.OrderStatusAccepted && order.Status != entity.OrderStatusPending {
		return nil, fmt.Errorf("заказ уже не ожидает принятия")
	}

	if cmd.Status == entity.OrderStatusPending {
		if order.DriverID == nil || cmd.DriverID == nil || *order.DriverID != *cmd.DriverID {
			return nil, fmt.Errorf("только назначенный водитель может отменить принятие заказа")
		}
		if order.Status != entity.OrderStatusAccepted && order.Status != entity.OrderStatusInProgress {
			return nil, fmt.Errorf("заказ не может быть возвращён в ожидание из статуса %s", order.Status)
		}
		if err := h.repo.DiscardDriverFromOrder(ctx, cmd.OrderID, *cmd.DriverID); err != nil {
			return nil, fmt.Errorf("ошибка при отмене принятия заказа: %w", err)
		}
		order, err = h.repo.GetOrderByID(ctx, cmd.OrderID)
		if err != nil {
			return nil, err
		}
		h.notifyDriverCancel(order, previousStatus, cmd.CancellationReason)
		return order, nil
	}

	if cmd.Status != entity.OrderStatusAccepted && order.DriverID != nil && cmd.DriverID != nil && *order.DriverID != *cmd.DriverID {
		return nil, fmt.Errorf("только назначенный водитель может изменить статус заказа")
	}

	if err := order.TransitionStatus(cmd.Status); err != nil {
		return nil, err
	}

	if cmd.Status == entity.OrderStatusAccepted {
		if cmd.DriverID == nil {
			return nil, fmt.Errorf("driver_id обязателен при принятии заказа")
		}
		order.DriverID = cmd.DriverID
		if err := h.repo.SetOrderDriver(ctx, cmd.OrderID, *cmd.DriverID); err != nil {
			return nil, fmt.Errorf("ошибка назначения водителя: %w", err)
		}
		h.sendNotifications(order, previousStatus)
		return order, nil
	}

	if cmd.Status == entity.OrderStatusCancelled {
		order.CancellationReason = cmd.CancellationReason
	}

	err = h.repo.UpdateOrderStatus(ctx, cmd.OrderID, cmd.Status, order.AcceptedAt, order.PickedUpAt, order.CompletedAt, order.CancelledAt, order.CancellationReason)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления статуса: %w", err)
	}

	h.sendNotifications(order, previousStatus)

	return order, nil
}

func (h *UpdateOrderStatusHandler) notifyDriverCancel(order *entity.Order, previousStatus entity.OrderStatus, reason *string) {
	if h.notifSvc == nil {
		return
	}
	go func() {
		ctx := context.Background()
		reasonText := "водитель не сможет выполнить заказ"
		if reason != nil && *reason != "" {
			reasonText = *reason
		}
		cancelMsg := fmt.Sprintf("Заказ #%d возвращён в поиск: %s", order.ID, reasonText)
		h.send(ctx, order.CustomerID, "Водитель отменил заказ", cancelMsg)
	}()
}

func (h *UpdateOrderStatusHandler) sendNotifications(order *entity.Order, previousStatus entity.OrderStatus) {
	if h.notifSvc == nil {
		return
	}

	go h.notifyStatusChange(order, previousStatus)
}

func (h *UpdateOrderStatusHandler) notifyStatusChange(order *entity.Order, previousStatus entity.OrderStatus) {
	ctx := context.Background()

	switch order.Status {
	case entity.OrderStatusAccepted:
		h.send(ctx, order.CustomerID,
			"Водитель назначен",
			fmt.Sprintf("Водитель принял ваш заказ #%d", order.ID))
		if order.DriverID != nil {
			h.send(ctx, *order.DriverID,
				"Новый заказ",
				fmt.Sprintf("Вы приняли заказ #%d", order.ID))
		}

	case entity.OrderStatusInProgress:
		h.send(ctx, order.CustomerID,
			"Водитель в пути",
			fmt.Sprintf("Водитель выехал к вам по заказу #%d", order.ID))

	case entity.OrderStatusCompleted:
		h.send(ctx, order.CustomerID,
			"Заказ завершён",
			fmt.Sprintf("Заказ #%d выполнен", order.ID))
		if order.DriverID != nil {
			h.send(ctx, *order.DriverID,
				"Заказ завершён",
				fmt.Sprintf("Заказ #%d успешно выполнен", order.ID))
		}

	case entity.OrderStatusCancelled:
		reason := ""
		if order.CancellationReason != nil {
			reason = *order.CancellationReason
		}
		cancelMsg := fmt.Sprintf("Заказ #%d отменён", order.ID)
		if reason != "" {
			cancelMsg = fmt.Sprintf("Заказ #%d отменён: %s", order.ID, reason)
		}
		h.send(ctx, order.CustomerID, "Заказ отменён", cancelMsg)
		if order.DriverID != nil {
			h.send(ctx, *order.DriverID, "Заказ отменён", cancelMsg)
		}
	}
}

func (h *UpdateOrderStatusHandler) send(ctx context.Context, userID int64, title, body string) {
	err := h.notifSvc.SendToUser(ctx, userID, notification.NotificationPayload{
		Title: title,
		Body:  body,
	})
	if err != nil {
		log.Printf("failed to notify user %d: %v", userID, err)
	}
}
