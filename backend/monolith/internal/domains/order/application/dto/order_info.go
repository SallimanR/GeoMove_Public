package dto

import geoDTO "monolith/internal/domains/geolocation/application/dto"

// Universal order info
type OrderInfo struct {
	Location geoDTO.Location
}
