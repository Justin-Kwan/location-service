package mocks

import (
	"time"
)

type TrackedItem struct {
	Location  Location
	ID        string
	CreatedAt int64
	UpdatedAt int64
}

func (c *TrackedItem) SetLocation(lon, lat float64) {
	c.Location.Lon = lon
	c.Location.Lat = lat
}

func (c *TrackedItem) SetCreatedAt() {
	c.CreatedAt = time.Now().UnixNano()
}

func (c *TrackedItem) SetUpdatedAt() {
	c.UpdatedAt = time.Now().UnixNano()
}

func (c *TrackedItem) GetID() string {
	return c.ID
}

func (c *TrackedItem) GetLon() float64 {
	return c.Location.Lon
}

func (c *TrackedItem) GetLat() float64 {
	return c.Location.Lat
}
