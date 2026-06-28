package dto

type LocationRaw [2]float32

type Location struct {
	Latitude  float32
	Longitude float32
}

type Coordinate uint8

const (
	Latitude  Coordinate = 0
	Longitude Coordinate = 1
)

