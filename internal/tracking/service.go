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

	// figure out how to stop
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for _ = range ticker.C {
			Idto, msgReceived := s.drain.Read()

			if msgReceived == false {
				continue
			}

			dto := Idto.(types.TrackCourierDTO)
			c.SetLocation(dto.Location.Lon, dto.Location.Lat)
			c.SetSpeed(dto.Speed)
			c.SetRadius(dto.Radius)
			c.SetUpdatedAt()

			s.upsertCourier(c)
			s.logCourierUpdate(c)
		}
	}()

	return nil
}

func (s *Service) logCourierUpdate(c *Courier) {
	log.Printf("-----------------------------------------------")
	log.Printf("Time: %v", time.Now().Format(time.RFC3339))
	log.Printf("New courier:")
	log.Printf("id: %s", c.GetID())
	log.Printf("coord: (%f, %f)", c.GetLon(), c.GetLat())
	log.Printf("speed: %f", c.Speed)
	log.Printf("radius: %f", c.Radius)
	log.Printf("created at: %d", c.CreatedAt)
	log.Printf("updated at: %d", c.UpdatedAt)
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

////////////////////// order functions ///////////////////////

// test!
// from message queue
func (s *Service) AddNewOrder(location map[string]float64, id string) error {
	o := NewOrder(id)

	err := s.orderStore.AddNew(o)
	return err
}

// test!
func (s *Service) DeleteOrder(id string) error {
	return s.orderStore.Delete(id)
}

// test!
func (s *Service) GetAllNearbyOrders(coord map[string]float64, radius float64) ([]string, error) {
	return s.orderStore.GetAllNearby(coord, radius)
}

// test!
func (s *Service) GetAllNearbyUnmatchedOrders(coord map[string]float64, radius float64) ([]string, error) {
	return s.orderStore.GetAllNearbyUnmatched(coord, radius)
}
