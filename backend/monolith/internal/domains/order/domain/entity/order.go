package entity

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderStatusForming    OrderStatus = "forming"
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAccepted   OrderStatus = "accepted"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID                   int64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CustomerID           int64
	DriverID             *int64
	FromLat              float32
	FromLon              float32
	FromAddress          string
	ToLat                float32
	ToLon                float32
	ToAddress            string
	TotalDistanceMeters  *int32
	HowManyWheelsBlocked int16
	PriceRubles          *int32
	Status               OrderStatus
	AcceptedAt           *time.Time
	PickedUpAt           *time.Time
	CompletedAt          *time.Time
	CancelledAt          *time.Time
	CancellationReason   *string
	CarWeightKg          int32
	CarLengthMeters      float32
	CarType              string
	CarName              string
	CarPhotoUrl          *string
	CustomerMessage      *string
}

type NewOrderOptions struct {
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

func NewOrder(opts NewOrderOptions) (*Order, error) {
	if opts.CustomerID <= 0 {
		return nil, fmt.Errorf("несуществеющий пользователь")
	}
	if opts.FromLat < -90 || opts.FromLat > 90 || opts.FromLon < -180 || opts.FromLon > 180 {
		return nil, fmt.Errorf("точка отправки находится за пределами карты")
	}
	if opts.ToLat < -90 || opts.ToLat > 90 || opts.ToLon < -180 || opts.ToLon > 180 {
		return nil, fmt.Errorf("точка прибытия находится за пределами карты")
	}
	if opts.HowManyWheelsBlocked <= 0 {
		return nil, fmt.Errorf("how_many_wheels_blocked must be positive")
	}

	now := time.Now()
	return &Order{
		CustomerID:           opts.CustomerID,
		FromLat:              opts.FromLat,
		FromLon:              opts.FromLon,
		FromAddress:          opts.FromAddress,
		ToLat:                opts.ToLat,
		ToLon:                opts.ToLon,
		ToAddress:            opts.ToAddress,
		HowManyWheelsBlocked: opts.HowManyWheelsBlocked,
		TotalDistanceMeters:  opts.TotalDistanceMeters,
		PriceRubles:          opts.PriceRubles,
		CarWeightKg:          opts.CarWeightKg,
		CarLengthMeters:      opts.CarLengthMeters,
		CarType:              opts.CarType,
		CarName:              opts.CarName,
		CarPhotoUrl:          opts.CarPhotoUrl,
		CustomerMessage:      opts.CustomerMessage,
		Status:               OrderStatusForming,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

var validTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusForming:    {OrderStatusPending, OrderStatusCancelled},
	OrderStatusPending:    {OrderStatusAccepted, OrderStatusCancelled},
	OrderStatusAccepted:   {OrderStatusInProgress},
	OrderStatusInProgress: {OrderStatusCompleted},
	OrderStatusCompleted:  {},
	OrderStatusCancelled:  {},
}

func (o *Order) TransitionStatus(newStatus OrderStatus) error {
	allowed, ok := validTransitions[o.Status]
	if !ok {
		return fmt.Errorf("unknown current status: %s", o.Status)
	}

	for _, s := range allowed {
		if s == newStatus {
			now := time.Now()
			o.Status = newStatus
			o.UpdatedAt = now

			switch newStatus {
			case OrderStatusAccepted:
				o.AcceptedAt = &now
			case OrderStatusInProgress:
				o.PickedUpAt = &now
			case OrderStatusCompleted:
				o.CompletedAt = &now
			case OrderStatusCancelled:
				o.CancelledAt = &now
			}
			return nil
		}
	}

	return fmt.Errorf("cannot transition from %s to %s", o.Status, newStatus)
}
