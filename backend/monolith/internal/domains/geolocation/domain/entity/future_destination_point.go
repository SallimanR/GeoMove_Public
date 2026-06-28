package entity

import "time"

type FutureDestinationPoint struct {
	Latitude          float32
	Longitude         float32
	TimeToArrive      time.Time
	TimeRangeToArrive time.Time
}

type DestinationFactory struct{}

func NewDestinationPoint() {}

func (fdp *FutureDestinationPoint) Validate() {}
