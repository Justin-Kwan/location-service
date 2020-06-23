package tracking

import (
	"time"
)

type Order struct {
	Location  Location
	ID        string
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

func NewOrder(id string) *Order {
	o := &Order{ID: id}
	o.SetCreatedAt()
	o.SetUpdatedAt()
	return o
}

func (o *Order) SetLocation(lon, lat float64) {
	o.Location.Lon = lon
	o.Location.Lat = lat
}

func (o *Order) SetCreatedAt() {
	o.CreatedAt = time.Now().UnixNano()
}

func (o *Order) SetUpdatedAt() {
	o.UpdatedAt = time.Now().UnixNano()
}

func (o *Order) GetID() string {
	return o.ID
}

func (o *Order) GetLon() float64 {
	return o.Location.Lon
}

func (o *Order) GetLat() float64 {
	return o.Location.Lat
}
