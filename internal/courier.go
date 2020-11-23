package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Courier struct {
	*TrackedItem
	Speed float64 `json:"speed"`
}

func NewCourier(id string) *Courier {
	return &Courier{
		TrackedItem: NewTrackedItem(id),
	}
}

func (c *Courier) SetSpeed(s float64) {
	c.Speed = s
}

func (c *Courier) GetSpeed() float64 {
	return c.Speed
}

func (c *Courier) Validate() error {
	if err := c.TrackedItem.Validate(); err != nil {
		return err
	}

	return validation.ValidateStruct(c,
		validation.Field(&c.Speed, validation.Required),
	)
}
