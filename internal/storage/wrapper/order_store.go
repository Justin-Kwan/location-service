package wrapper

import "location-service/internal"

type OrderStore struct {
	itemStore ItemStorer
}

func NewOrderStore(itemStore ItemStorer) *OrderStore {
	return &OrderStore{
		itemStore: itemStore,
	}
}

func (os *OrderStore) FindAllNearbyOrderIDs(coord *internal.Location, radius float64) ([]string, error) {
	orderIDs, err := os.itemStore.findAllNearbyItemIDs(coord, radius)
	return orderIDs, err
}

func (os *OrderStore) FindAllNearbyUnmatchedOrderIDs(coord *internal.Location, radius float64) ([]string, error) {
	// TODO: https://swiperun.atlassian.net/browse/SR-2

	return []string{"mock_order1"}, nil
}

func (os *OrderStore) AddNewOrder(o *internal.Order) error {
	// TODO: https://swiperun.atlassian.net/browse/SR-4

	return nil
}

func (os *OrderStore) GetOrder(id string) (*internal.Order, error) {
	// TODO: https://swiperun.atlassian.net/browse/SR-8

	return internal.NewOrder("mock_order_id"), nil
}

func (os *OrderStore) DeleteOrder(id string) error {
	// TODO: https://swiperun.atlassian.net/browse/SR-7

	return nil
}
