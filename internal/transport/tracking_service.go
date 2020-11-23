package transport

//go:generate mockgen -destination=../mocks/mock_tracking_service.go -package=mocks . TrackingService

import (
	"location-service/internal"
)

// TrackingService describes behaviour of domain services.
type TrackingService interface {
	TrackCourier(id string) error
	DeleteCourier(id string) error
	GetAllNearbyCouriers(coord *internal.Location, radius float64) ([]string, error)
	FindAllNearbyOrderIDs(coord *internal.Location, radius float64) ([]string, error)
	AddNewOrder(o *internal.Order) error
	DeleteOrder(id string) error
	// GetAllNearbyUnmatchedOrders(r *tracking.FindAllNearbyRequest) ([]string, error)
}
