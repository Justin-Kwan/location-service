package internal

import validation "github.com/go-ozzo/ozzo-validation"

const (
	maxLon = 180.0
	minLon = -180.0
	maxLat = 90.0
	minLat = -90.0
)

type Location struct {
	Lon float64
	Lat float64
}

func NewLocation(lon, lat float64) Location {
	return Location{
		Lon: lon,
		Lat: lat,
	}
}

func (l Location) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Lon, validation.Required, validation.Max(maxLon)),
		validation.Field(&l.Lon, validation.Required, validation.Min(minLon)),
		validation.Field(&l.Lat, validation.Required, validation.Max(maxLat)),
		validation.Field(&l.Lat, validation.Required, validation.Min(minLat)),
	)
}
