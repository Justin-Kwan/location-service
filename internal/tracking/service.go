package tracking

import (
	"log"
	"time"

	"location-service/internal"
	"location-service/internal/types"

	"github.com/google/uuid"
)

var newUUID = uuid.NewRandom

type Service struct {
	courierStore CourierStorer
	orderStore   OrderStorer
	drain        types.Drain
}

func NewService(cs CourierStorer, os OrderStorer, dn types.Drain) *Service {
	return &Service{
		courierStore: cs,
		orderStore:   os,
		drain:        dn,
	}
}

func (s *Service) GetCourier(id string) (*internal.Courier, error) {
	c, err := s.courierStore.GetCourier(id)
	return c, err
}

func (s *Service) TrackCourier(id string) error {
	c := internal.NewCourier(id)

	// figure out how to stop
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for range ticker.C {
			Idto, msgReceived := s.drain.Read()

			if msgReceived == false {
				continue
			}

			dto := Idto.(types.TrackCourierDTO)
			c.SetLocation(dto.Location.Lon, dto.Location.Lat)
			c.SetSpeed(dto.Speed)
			c.SetSearchRadius(dto.Radius)
			c.SetUpdatedAt()

			s.upsertCourier(c)
			s.logCourierUpdate(c)
		}
	}()

	return nil
}

func (s *Service) logCourierUpdate(c *internal.Courier) {
	log.Printf("-----------------------------------------------")
	log.Printf("Time: %v", time.Now().Format(time.RFC3339))
	log.Printf("New courier:")
	log.Printf("id: %s", c.GetID())
	log.Printf("coord: (%f, %f)", c.GetLon(), c.GetLat())
	log.Printf("speed: %f", c.Speed)
	log.Printf("radius: %f", c.GetSearchRadius())
	log.Printf("created at: %d", c.CreatedAt)
	log.Printf("updated at: %d", c.UpdatedAt)
	log.Printf("-----------------------------------------------")
}

// move to db layer
func (s *Service) upsertCourier(c *internal.Courier) error {
	// err := s.courierStore.Update(c)

	// if err != nil && err == types.ErrKeyNotFound {
	return s.courierStore.UpsertCourier(c)
	// }

	// return err
}

// DeleteCourier ... (comes from message queue)
func (s *Service) DeleteCourier(id string) error {
	return s.courierStore.DeleteCourier(id)
}

// GetAllNearbyCouriers ...
func (s *Service) GetAllNearbyCouriers(coord *internal.Location, radius float64) ([]string, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	id := uuid.String()

	t := internal.NewTrackedItem(id)
	t.SetLocation(coord.Lon, coord.Lat)

	courierIDs, err := s.courierStore.FindAllNearbyCourierIDs(t, radius)
	return courierIDs, err
}

////////////////////// order functions ///////////////////////

// AddNewOrder ... (from message queue)
func (s *Service) AddNewOrder(o *internal.Order) error {
	// dupe check, validation

	err := s.orderStore.AddNewOrder(o)
	return err
}

// DeleteOrder ...
func (s *Service) DeleteOrder(id string) error {
	return s.orderStore.DeleteOrder(id)
}

// FindAllNearbyOrderIDs returns all order ids nearby to
func (s *Service) FindAllNearbyOrderIDs(coord *internal.Location, radius float64) ([]string, error) {
	if err := coord.Validate(); err != nil {
		return nil, err
	}

	// validate radius!

	orderIDs, err := s.orderStore.FindAllNearbyOrderIDs(coord, radius)
	return orderIDs, err
}

// FindAllNearbyOrderIDsToDriver ...
func (s *Service) FindAllNearbyOrderIDsToDriver(courierID string) ([]string, error) {
	c, err := s.GetCourier(courierID)
	if err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	coord := c.GetLocation()
	searchRadius := c.GetSearchRadius()

	log.Println("<<<<<<<<<<<<<<< LOCATION SERVICE >>>>>>>>>>")
	log.Println("COORD REQUESTED", coord)
	log.Println("SEARCH RADIUS", searchRadius)

	ids, err := s.orderStore.FindAllNearbyUnmatchedOrderIDs(&coord, searchRadius)
	return ids, err
}

// // GetAllNearbyUnmatchedOrders ...
// func (s *Service) FindAllNearbyUnmatchedOrders(t) ([]string, error) {
// 	t := internal.NewTrackedItem(r.ID)
// 	t.SetLocation(r.Coord.Lon, r.Coord.Lat)

// 	return s.orderStore.FindAllNearbyUnmatchedOrderIDs(t, r.Radius)
// }
