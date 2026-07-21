package command

import (
	"context"
	"fmt"
	"log"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
	"monolith/internal/notification"
)

type CreateOrderCommand struct {
	CustomerID           int64
	FromLat              float32
	FromLon              float32
	FromAddress          string
	ToLat                float32
	ToLon                float32
	ToAddress            string
	HowManyWheelsBlocked int16
	TotalDistanceMeters  *int32
	PriceRubles          *int32
	CarWeightKg          int32
	CarLengthMeters      float32
	CarType              string
	CarName              string
	CarPhotoUrl          *string
	CustomerMessage      *string
}

type CreateOrderHandler struct {
	repo     repository.OrderRepository
	notifSvc *notification.Service
}

func NewCreateOrderHandler(repo repository.OrderRepository, notifSvc *notification.Service) *CreateOrderHandler {
	return &CreateOrderHandler{repo: repo, notifSvc: notifSvc}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error) {
	existing, err := h.repo.ListOrdersByCustomer(ctx, cmd.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказов: %w", err)
	}
	for _, o := range existing {
		if o.Status == entity.OrderStatusForming || o.Status == entity.OrderStatusPending {
			return nil, fmt.Errorf("у вас уже есть активный заказ #%d", o.ID)
		}
	}

	order, err := entity.NewOrder(entity.NewOrderOptions{
		CustomerID:           cmd.CustomerID,
		FromLat:              cmd.FromLat,
		FromLon:              cmd.FromLon,
		FromAddress:          cmd.FromAddress,
		ToLat:                cmd.ToLat,
		ToLon:                cmd.ToLon,
		ToAddress:            cmd.ToAddress,
		HowManyWheelsBlocked: cmd.HowManyWheelsBlocked,
		TotalDistanceMeters:  cmd.TotalDistanceMeters,
		PriceRubles:          cmd.PriceRubles,
		CarWeightKg:          cmd.CarWeightKg,
		CarLengthMeters:      cmd.CarLengthMeters,
		CarType:              cmd.CarType,
		CarName:              cmd.CarName,
		CarPhotoUrl:          cmd.CarPhotoUrl,
		CustomerMessage:      cmd.CustomerMessage,
	})
	if err != nil {
		return nil, err
	}

	id, err := h.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	order.ID = id

	if h.notifSvc != nil {
		go func() {
			err := h.notifSvc.SendToUser(context.Background(), cmd.CustomerID, notification.NotificationPayload{
				Title: "Заказ создан",
				Body:  fmt.Sprintf("Ваш заказ #%d создан и ожидает водителя", id),
			})
			if err != nil {
				log.Printf("failed to notify customer %d about order %d: %v", cmd.CustomerID, id, err)
			}
		}()
	}

	return order, nil
}
