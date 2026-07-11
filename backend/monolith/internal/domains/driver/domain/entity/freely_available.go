package entity

import "time"

type LocationWithAddress struct {
	Lat     float32
	Lon     float32
	Address string
}

type FreelyAvailable struct {
	UserID       int64
	FromDate     time.Time
	ToDate       time.Time
	FromLocation LocationWithAddress
	ToLocations  []LocationWithAddress
	EnRouteOrder *bool
	TariffPerKm  *float32
}

type FreelyAvailableDriver struct {
	UserID       int64
	FromDate     time.Time
	ToDate       time.Time
	FromLocation LocationWithAddress
	EnRouteOrder *bool
	TariffPerKm  *float32
	Name         string
	Rating       *float32
	ProfileImage *string
	Distance     float32
}

func NewFreelyAvailable(userID int64, fromDate, toDate time.Time, fromLocation LocationWithAddress, toLocations []LocationWithAddress, enRouteOrder *bool, tariffPerKm *float32) *FreelyAvailable {
	if toLocations == nil {
		toLocations = []LocationWithAddress{}
	}
	return &FreelyAvailable{
		UserID:       userID,
		FromDate:     fromDate,
		ToDate:       toDate,
		FromLocation: fromLocation,
		ToLocations:  toLocations,
		EnRouteOrder: enRouteOrder,
		TariffPerKm:  tariffPerKm,
	}
}
