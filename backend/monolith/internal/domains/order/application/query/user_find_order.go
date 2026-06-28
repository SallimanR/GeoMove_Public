package query

import (
	"time"

	geoDTO "monolith/internal/domains/geolocation/application/dto"
	geoRepo "monolith/internal/domains/geolocation/domain/repository"
)

type FindOrderQuery struct {
	UserLocation       geoDTO.Location
	OrderScheduledTime *time.Time
	// OrderInfo
}

type FindOrderHandler struct {
	geoRepo geoRepo.GeolocationRepository
}

func NewFindOrderHandler(geoRepo geoRepo.GeolocationRepository) *FindOrderHandler {
	return &FindOrderHandler{
		geoRepo: geoRepo,
	}
}

func (h *FindOrderHandler) Handle() {}
