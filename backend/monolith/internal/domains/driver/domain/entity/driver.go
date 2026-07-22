package entity

import (
	"context"
	"fmt"
	"time"

	"monolith/pkg/geo"
)

type DriverID uint32

type Location struct {
	Lat float32
	Lon float32
}

type Driver struct {
	UserID              int64
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	Name                string
	ProfileImage        *string
	WorkStarts          *time.Time
	WorkEnds            *time.Time
	IsAvailable         bool
	LastSeen            time.Time
	Location            Location
	Rating             float32
	Phone              string
	MaxCarWeightKg     int32
	MaxCarLengthMeters float32
	Address            string
}

type DriverOptions struct {
	UserID     int64
	Name       string
	Phone      *string
	Location   Location
	WorkStarts *time.Time
	WorkEnds   *time.Time
	Address    string
}

func NewDriver(options DriverOptions) (*Driver, error) {
	now := time.Now()
	phone := ""
	if options.Phone != nil {
		phone = *options.Phone
	}
	driver := &Driver{
		UserID:      options.UserID,
		Name:        options.Name,
		CreatedAt:   now,
		LastSeen:    now,
		IsAvailable: true,
		WorkStarts:  options.WorkStarts,
		WorkEnds:    options.WorkEnds,
		Location:    options.Location,
		Rating:      0,
		Phone:       phone,
		Address:     options.Address,
	}

	err := driver.Validate()
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (d *Driver) Validate() error {
	if d.Location.Lat > 90 || d.Location.Lat < -90 || d.Location.Lon > 180 || d.Location.Lon < -180 {
		return fmt.Errorf("coordinates are out of range")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	countryCode, err := geo.ResolveCountry(ctx, d.Location.Lat, d.Location.Lon)
	if err != nil {
		return err
	}
	const countryCodeValid = "RU"
	if countryCode != countryCodeValid {
		return fmt.Errorf("invalid country")
	}

	return nil
}

func (d *Driver) UpdateWorkHours(start, end *time.Time) {
	d.WorkStarts = start
	d.WorkEnds = end
	d.UpdatedAt = timeNowPtr()
}

func (d *Driver) UpdateLocation(lat, lon float32) error {
	d.Location.Lat = lat
	d.Location.Lon = lon
	err := d.Validate()
	if err != nil {
		return err
	}
	return nil
}

func timeNowPtr() *time.Time {
	t := time.Now()
	return &t
}
