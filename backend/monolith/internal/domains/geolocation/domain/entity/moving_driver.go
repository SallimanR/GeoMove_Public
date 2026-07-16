package entity

import "time"

type MovingDriver struct {
	DriverID   int64
	Latitude   float32
	Longitude  float32
	TravelTime time.Time
	PathMeters int32
}

type MovingDriverWithPoints struct {
	MovingDriver
	Points [][2]float32
}

func NewMovingDriver(driverID int64, gpsPoints []LocationPoint, travelTime time.Time, pathMeters uint32) *MovingDriver {
	pointsNumber := len(gpsPoints)
	if pointsNumber < 2 {
		return nil
	}

	lastPoint := pointsNumber - 1
	return &MovingDriver{
		DriverID:   driverID,
		Latitude:   gpsPoints[lastPoint].Latitude,
		Longitude:  gpsPoints[lastPoint].Longitude,
		TravelTime: travelTime,
		PathMeters: int32(pathMeters),
	}
}
