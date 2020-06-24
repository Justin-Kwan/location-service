package tracking

type Courier struct {
	TrackedItem
	Speed       float64 `json:"speed"`
	Radius      float64 `json:"radius"`
}

func NewCourier(id string) *Courier {
	return &Courier{
		TrackedItem: NewTrackedItem(id),
	}
}

func (c *Courier) SetRadius(r float64) {
	c.Radius = r
}

func (c *Courier) SetSpeed(s float64) {
	c.Speed = s
}
