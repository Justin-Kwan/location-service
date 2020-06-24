package tracking

type Order struct {
	TrackedItem
}

func NewOrder(id string) *Order {
	return &Order{
		TrackedItem: NewTrackedItem(id),
	}
}
