package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestTrackedItems() []TrackedItem {
	return []TrackedItem{
		TrackedItem{
			Coord:     Location{Lon: 67.1349, Lat: 78.45314},
			ID:        "test_tracked_item_id1",
			CreatedAt: 123456789009,
			UpdatedAt: 987654321001,
		},
		TrackedItem{
			Coord:     Location{Lon: 43.9876567888854435656, Lat: 65.7654567876},
			ID:        "test_tracked_item_id2",
			CreatedAt: 123456787419,
			UpdatedAt: 987892321923,
		},
		TrackedItem{
			Coord:     Location{Lon: -180, Lat: -85.05112878},
			ID:        "test_tracked_item_id3",
			CreatedAt: 123456789009,
			UpdatedAt: 987612432101,
		},
	}
}

func TestNewTrackedItem_NormalCases(t *testing.T) {
	setupMinUnixNanoTime()
	ti := NewTrackedItem("test_id")

	assert.Equal(t, ti.GetID(), "test_id")
	assert.GreaterOrEqual(t, ti.CreatedAt, MinUnixNanoTime)
	assert.GreaterOrEqual(t, ti.UpdatedAt, MinUnixNanoTime)
}

func TestSetLocation_NormalCases(t *testing.T) {
	tis := getTestTrackedItems()

	for _, ti := range tis {
		// function under test
		ti.SetLocation(TestLon, TestLat)

		// assert correct location set
		assert.Equal(t, TestLon, ti.GetLon())
		assert.Equal(t, TestLat, ti.GetLat())
	}
}

func TestSetCreatedAt_NormalCases(t *testing.T) {
	setupMinUnixNanoTime()
	tis := getTestTrackedItems()

	for _, ti := range tis {
		// function under test
		ti.SetCreatedAt()

		// assert correct location set
		assert.GreaterOrEqual(t, ti.CreatedAt, MinUnixNanoTime)
	}
}

func TestSetUpdatedAt_NormalCases(t *testing.T) {
	setupMinUnixNanoTime()
	tis := getTestTrackedItems()

	for _, ti := range tis {
		// function under test
		ti.SetUpdatedAt()

		// assert correct location set
		assert.GreaterOrEqual(t, ti.UpdatedAt, MinUnixNanoTime)
	}
}

func TestGetID_NormalCases(t *testing.T) {
	tis := getTestTrackedItems()

	for _, ti := range tis {
		// function under test
		assert.Equal(t, ti.GetID(), ti.ID)
	}
}
