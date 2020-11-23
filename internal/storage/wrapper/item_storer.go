package wrapper

//go:generate mockgen -destination=../../mocks/mock_item_storer.go -package=mocks . ItemStorer

import "location-service/internal"

type ItemStorer interface {
	addNewItem(t *internal.TrackedItem) error
	getItem(id string) (*internal.TrackedItem, error)
	findAllNearbyItemIDs(coord *internal.Location, radius float64) ([]string, error)
	getUnmatchedNearby(coord map[string]float64, radius float64) ([]string, error)
	update(t *internal.TrackedItem) error
	delete(id string) error
	setUnmatched(id string) error
	setMatched(id string) error
}
