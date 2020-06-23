package mock

import (
  "time"
)

type Courier struct {
  Location  Location
	ID        string
	Speed     float64
	Radius    float64
	CreatedAt int64
	UpdatedAt int64
}

func (c *Courier) SetLocation(lon, lat float64) {
	c.Location.Lon = lon
	c.Location.Lat = lat
}

func (c *Courier) SetRadius(r float64) {
  c.Radius = r
}

func (c *Courier) SetSpeed(s float64) {
  c.Speed = s
}

func (c *Courier) SetCreatedAt() {
	c.CreatedAt = time.Now().UnixNano()
}

func (c *Courier) SetUpdatedAt() {
	c.UpdatedAt = time.Now().UnixNano()
}

func (c *Courier) GetID() string {
	return c.ID
}

func (c *Courier) GetLon() float64 {
	return c.Location.Lon
}

func (c *Courier) GetLat() float64 {
	return c.Location.Lat
}
