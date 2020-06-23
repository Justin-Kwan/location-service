package tracking

import (
	"testing"

	// "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func getTestOrders() []Order {
	return []Order{
		Order{
			Location:  Location{Lon: 67.1349, Lat: 78.45314},
			ID:        "test_order_id1",
			CreatedAt: 123456789009,
			UpdatedAt: 987654321001,
		},
    Order{
			Location:  Location{Lon: 43.9876567888854435656, Lat: 65.7654567876},
			ID:        "test_order_id2",
			CreatedAt: 123456787419,
			UpdatedAt: 987892321923,
		},
    Order{
			Location:  Location{Lon: -180, Lat: -85.05112878},
			ID:        "test_order_id3",
			CreatedAt: 123456789009,
			UpdatedAt: 987612432101,
		},
	}
}

func TestOrderSetLocation(t *testing.T) {
	tos := getTestOrders()

	for _, to := range tos {
		// function under test
		to.SetLocation(TestLon, TestLat)

		// assert correct location set
		assert.Equal(t, TestLon, to.Location.Lon)
		assert.Equal(t, TestLat, to.Location.Lat)
	}
}

func TestOrderSetCreatedAt(t *testing.T) {
	setupMinUnixNanoTime()
	tos := getTestOrders()

	for _, to := range tos {
		// function under test
		to.SetCreatedAt()

		// assert correct location set
		assert.GreaterOrEqual(t, to.CreatedAt, MinUnixNanoTime)
	}
}

func TestOrderSetUpdatedAt(t *testing.T) {
	setupMinUnixNanoTime()
	tos := getTestOrders()

	for _, to := range tos {
		// function under test
		to.SetUpdatedAt()

		// assert correct location set
		assert.GreaterOrEqual(t, to.UpdatedAt, MinUnixNanoTime)
	}
}
