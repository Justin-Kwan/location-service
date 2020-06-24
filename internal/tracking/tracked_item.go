package tracking

import (
	"time"
)

type TrackedItem struct {
	Coord     Location
	ID        string
	CreatedAt int64
	UpdatedAt int64
}

type Location struct {
	Lon float64
	Lat float64
}

func NewTrackedItem(id string) TrackedItem {
	return TrackedItem{
		ID:        id,
		CreatedAt: time.Now().UnixNano(),
		UpdatedAt: time.Now().UnixNano(),
	}
}

func (t *TrackedItem) SetLocation(lon, lat float64) {
	t.Coord.Lon = lon
	t.Coord.Lat = lat
}

func (t *TrackedItem) SetCreatedAt() {
	t.CreatedAt = time.Now().UnixNano()
}

func (t *TrackedItem) SetUpdatedAt() {
	t.UpdatedAt = time.Now().UnixNano()
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
