package internal

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// TODO: move to config file
	maxSearchRadiusKm = 50.0
	minSearchRadiusKm = 1.0
)

var (
	UnixNanoNow = time.Time.UnixNano
	uuidRegex   = regexp.MustCompile("^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$")
)

// TrackedItem represents any entity within location service
// that can be located to a coordinate location.
type TrackedItem struct {
	Coord        Location
	SearchRadius float64
	ID           string
	CreatedAt    int64
	UpdatedAt    int64
}

func NewTrackedItem(id string) *TrackedItem {
	return &TrackedItem{
		ID:        id,
		CreatedAt: UnixNanoNow(time.Now()),
		UpdatedAt: UnixNanoNow(time.Now()),
	}
}

func (t *TrackedItem) SetLocation(lon, lat float64) {
	t.Coord = NewLocation(lon, lat)
}

func (t *TrackedItem) SetSearchRadius(radius float64) {
	t.SearchRadius = radius
}

func (t *TrackedItem) SetCreatedAt() {
	t.CreatedAt = UnixNanoNow(time.Now())
}

func (t *TrackedItem) SetUpdatedAt() {
	t.UpdatedAt = UnixNanoNow(time.Now())
}

func (t *TrackedItem) GetID() string {
	return t.ID
}

func (t *TrackedItem) GetLon() float64 {
	return t.Coord.Lon
}

func (t *TrackedItem) GetLat() float64 {
	return t.Coord.Lat
}

func (t *TrackedItem) GetLocation() Location {
	return t.Coord
}

func (t *TrackedItem) GetSearchRadius() float64 {
	return t.SearchRadius
}

func (t *TrackedItem) Validate() error {
	if err := t.Coord.Validate(); err != nil {
		return err
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.SearchRadius, validation.Required, validation.Max(maxSearchRadiusKm)),
		validation.Field(&t.SearchRadius, validation.Required, validation.Min(minSearchRadiusKm)),
		validation.Field(&t.ID, validation.Required, validation.Match(uuidRegex)),
		validation.Field(&t.CreatedAt, validation.Required, validation.Max(UnixNanoNow(time.Now()))),
		validation.Field(&t.UpdatedAt, validation.Required, validation.Max(UnixNanoNow(time.Now()))),
	)
}
