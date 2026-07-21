package dto

import (
	"time"

	"monolith/internal/domains/order/domain/entity"
)

type OrderDTO struct {
	ID                   int64               `json:"id"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	CustomerID           int64               `json:"customer_id"`
	DriverID             *int64              `json:"driver_id,omitempty"`
	FromLat              float32             `json:"from_lat"`
	FromLon              float32             `json:"from_lon"`
	ToLat                float32             `json:"to_lat"`
	ToLon                float32             `json:"to_lon"`
	TotalDistanceMeters  *int32              `json:"total_distance_meters,omitempty"`
	HowManyWheelsBlocked int16               `json:"how_many_wheels_blocked"`
	PriceRubles          *int32              `json:"price_rubles,omitempty"`
	CarWeightKg          int32               `json:"car_weight_kg"`
	CarLengthMeters      float32             `json:"car_length_meters"`
	CarType              string              `json:"car_type"`
	CarName              string              `json:"car_name"`
	CarPhotoUrl          *string             `json:"car_photo_url,omitempty"`
	CustomerMessage      *string             `json:"customer_message,omitempty"`
	Status               entity.OrderStatus  `json:"status"`
	AcceptedAt           *time.Time          `json:"accepted_at,omitempty"`
	PickedUpAt           *time.Time          `json:"picked_up_at,omitempty"`
	CompletedAt          *time.Time          `json:"completed_at,omitempty"`
	CancelledAt          *time.Time          `json:"cancelled_at,omitempty"`
	CancellationReason   *string             `json:"cancellation_reason,omitempty"`
}

func FromEntity(o *entity.Order) OrderDTO {
	return OrderDTO{
		ID:                   o.ID,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
		CustomerID:           o.CustomerID,
		DriverID:             o.DriverID,
		FromLat:              o.FromLat,
		FromLon:              o.FromLon,
		ToLat:                o.ToLat,
		ToLon:                o.ToLon,
		TotalDistanceMeters:  o.TotalDistanceMeters,
		HowManyWheelsBlocked: o.HowManyWheelsBlocked,
		PriceRubles:          o.PriceRubles,
		CarWeightKg:          o.CarWeightKg,
		CarLengthMeters:      o.CarLengthMeters,
		CarType:              o.CarType,
		CarName:              o.CarName,
		CarPhotoUrl:          o.CarPhotoUrl,
		CustomerMessage:      o.CustomerMessage,
		Status:               o.Status,
		AcceptedAt:           o.AcceptedAt,
		PickedUpAt:           o.PickedUpAt,
		CompletedAt:          o.CompletedAt,
		CancelledAt:          o.CancelledAt,
		CancellationReason:   o.CancellationReason,
	}
}
