package wrapper

import "location-service/internal"

type CourierStore struct {
	itemStore ItemStorer
}

func NewCourierStore(itemStore ItemStorer) *CourierStore {
	return &CourierStore{
		itemStore: itemStore,
	}
}

func (cs *CourierStore) FindAllNearbyCourierIDs(t *internal.TrackedItem, radius float64) ([]string, error) {
	// TODO: https://swiperun.atlassian.net/browse/SR-6

	return []string{"mock_courier1"}, nil
}

func (cs *CourierStore) UpsertCourier(c *internal.Courier) error {
	// TODO: https://swiperun.atlassian.net/browse/SR-4

	return nil
}

func (cs *CourierStore) GetCourier(id string) (*internal.Courier, error) {
	// TODO: https://swiperun.atlassian.net/browse/SR-8

	return internal.NewCourier("mock_courier_id"), nil
}

func (cs *CourierStore) DeleteCourier(id string) error {
	// TODO: https://swiperun.atlassian.net/browse/SR-7

	return nil
}
