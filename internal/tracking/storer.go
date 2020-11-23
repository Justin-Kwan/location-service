package tracking

//go:generate mockgen -destination=../mocks/mock_courier_storer.go -package=mocks . CourierStorer
//go:generate mockgen -destination=../mocks/mock_order_storer.go -package=mocks . OrderStorer

import "location-service/internal"

type CourierStorer interface {
	FindAllNearbyCourierIDs(t *internal.TrackedItem, radius float64) ([]string, error)
	UpsertCourier(c *internal.Courier) error
	GetCourier(id string) (*internal.Courier, error)
	DeleteCourier(id string) error
}

type OrderStorer interface {
	FindAllNearbyOrderIDs(coord *internal.Location, radius float64) ([]string, error)
	FindAllNearbyUnmatchedOrderIDs(coord *internal.Location, radius float64) ([]string, error)
	AddNewOrder(o *internal.Order) error
	GetOrder(id string) (*internal.Order, error)
	DeleteOrder(id string) error
}
