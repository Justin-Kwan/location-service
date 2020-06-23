package tracking

import (
	"log"
	"time"

	"location-service/internal/types"
)

type Service struct {
	courierStore types.ItemStorer
	orderStore   types.ItemStorer
	drain        types.Drain
}

type TrackCourierDTO struct {
	Location struct {
		Lon float64
		Lat float64
	}
	Speed  float64
	Radius float64
}

func NewService(cs, os types.ItemStorer, dn types.Drain) *Service {
	return &Service{
		courierStore: cs,
		orderStore:   os,
		drain:        dn,
	}
}

// tested
func (s *Service) GetCourier(id string) (*Courier, error) {
	c := &Courier{}
	err := s.courierStore.Get(id, c)
	return c, err
}

// test!
func (s *Service) TrackCourier(id string) error {
	c := NewCourier(id)

	go func() {
		for {
			log.Printf("Reading from drain")
			dto := s.drain.Read().(TrackCourierDTO)

			c.SetLocation(dto.Location.Lon, dto.Location.Lat)
			c.SetSpeed(dto.Speed)
			c.SetRadius(dto.Radius)
			c.SetUpdatedAt()

			log.Printf("-----------------------------------------------")
			log.Printf("New courier:")
			log.Printf("id: %s", c.ID)
			log.Printf("coord: (%f, %f)", c.Location.Lon, c.Location.Lat)
			log.Printf("speed: %f", c.Speed)
			log.Printf("radius: %f", c.Radius)
			log.Printf("created at: %d", c.CreatedAt)
			log.Printf("updated at: %d", c.UpdatedAt)

			s.upsertCourier(c)
			time.Sleep(2 * time.Second)
		}
	}()

	return nil
}

// tested
func (s *Service) upsertCourier(c *Courier) error {
	err := s.courierStore.Update(c)

	if err != nil && err == types.ErrKeyNotFound {
		return s.courierStore.AddNew(c)
	}

	return err
}

// tested
// comes from message queue
func (s *Service) DeleteCourier(id string) error {
	return s.courierStore.Delete(id)
}

// better args?
// tested
func (s *Service) GetAllNearbyCouriers(coord map[string]float64, radius float64) ([]string, error) {
	return s.courierStore.GetAllNearby(coord, radius)
}

////////////////////// order functions ////////////////////

// test
// func (s *Service) AddNewOrder() error {
//
// }

// test!
func (s *Service) DeleteOrder(id string) error {
	return s.orderStore.Delete(id)
}

func (s *Service) GetAllNearbyUnmatchedOrders(coord map[string]float64, radius float64) ([]string, error) {
	return s.orderStore.GetAllNearbyUnmatched(coord, radius)
}

func (s *Service) GetAllNearbyOrders(coord map[string]float64, radius float64) ([]string, error) {
	return s.orderStore.GetAllNearby(coord, radius)
}
