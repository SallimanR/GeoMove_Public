package entity

import "monolith/internal/domains/geolocation/application/dto"

type LocationPoint struct {
	Latitude  float32
	Longitude float32
}

func ArrayToPoints(coords []dto.LocationRaw) *[]LocationPoint {
	points := make([]LocationPoint, len(coords))
	for i, coord := range coords {
		points[i].Longitude = coord[0]
		points[i].Latitude = coord[1]
	}
	return &points
}
