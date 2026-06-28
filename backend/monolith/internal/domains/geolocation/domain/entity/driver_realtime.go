package entity

import (
	"monolith/pkg/geo"
)

type DriverRealtime struct {
	DriverID     uint32
	Latitude     float32
	Longitude    float32
	Bearing      float32
	AverageSpeed float32
}

func NewDriverRealtime(driverID uint32, gpsPoints []LocationPoint, totalTime, totalDistance float32) *DriverRealtime {
	pointsNumber := len(gpsPoints)
	if pointsNumber < 2 {
		return nil
	}

	lastPoint := pointsNumber - 1
	return &DriverRealtime{
		DriverID:  driverID,
		Latitude:  gpsPoints[lastPoint].Latitude,
		Longitude: gpsPoints[lastPoint].Longitude,
		Bearing: geo.CalculateBearing(
			float64(gpsPoints[0].Latitude),
			float64(gpsPoints[0].Longitude),
			float64(gpsPoints[lastPoint].Latitude),
			float64(gpsPoints[lastPoint].Longitude),
		),
		AverageSpeed: totalDistance / totalTime,
	}
}

func (gr *DriverRealtime) Validate() {
	if gr.AverageSpeed < 5 {
	}
}
