package tracking

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TestLon = 54.4124
	TestLat = -45.41239

	TestRadius = 15.241512
	TestSpeed  = 5.0995
)

var (
	MinUnixNanoTime int64
)

func getTestCouriers() []Courier {
	return []Courier{
		Courier{
			TrackedItem: TrackedItem{
				Coord:     Location{Lon: 67.1349, Lat: 78.45314},
				ID:        "test_courier_id1",
				CreatedAt: 123456789009,
				UpdatedAt: 987654321001,
			},
			Speed:  45.12,
			Radius: 30,
		},
		Courier{
			TrackedItem: TrackedItem{
				Coord:     Location{Lon: 43.9876567888854435656, Lat: 65.7654567876},
				ID:        "test_courier_id2",
				CreatedAt: 987656,
				UpdatedAt: 8,
			},
			Speed:  0.0000000000001,
			Radius: 0.0000000000002,
		},
		Courier{
			TrackedItem: TrackedItem{
				Coord:     Location{Lon: -180, Lat: -85.05112878},
				ID:        "test_courier_id3",
				CreatedAt: 987656765678,
				UpdatedAt: 8765676567876,
			},
			Speed:  34.1231412,
			Radius: 10.14,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id4",
				Coord:     Location{Lon: 180, Lat: 85.05112878},
				CreatedAt: 987656765678,
				UpdatedAt: 8765676567876,
			},
			Speed:  34.1231412,
			Radius: 10.14,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id5",
				Coord:     Location{Lon: -180, Lat: 85.05112878},
				CreatedAt: 987656765678,
				UpdatedAt: 8765676567876,
			},
			Speed:  34.1231412,
			Radius: 10.14,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id6",
				Coord:     Location{Lon: 180, Lat: -85.05112878},
				CreatedAt: 987656765678,
				UpdatedAt: 8765676567876,
			},
			Speed:  34.1231412,
			Radius: 10.14,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id7",
				Coord:     Location{Lon: -79.661522, Lat: 43.458401},
				CreatedAt: 1591933701672,
				UpdatedAt: 1591933701672,
			},
			Speed:  50,
			Radius: 15,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id7.5",
				Coord:     Location{Lon: 0, Lat: 0},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			Speed:  0,
			Radius: 0,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id8",
				Coord:     Location{Lon: -79.661522, Lat: 43.458401},
				CreatedAt: 1991933701672,
				UpdatedAt: 1591933709672,
			},
			Speed:  50,
			Radius: 15,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id9",
				Coord:     Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 1591963701672,
				UpdatedAt: 1592933701672,
			},
			Speed:  50,
			Radius: 15,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id10",
				Coord:     Location{Lon: -80.481522, Lat: 43.328401},
				CreatedAt: 1591933701671,
				UpdatedAt: 1591933701672,
			},
			Speed:  20,
			Radius: 10,
		},
		Courier{
			TrackedItem: TrackedItem{
				ID:        "test_courier_id11",
				Coord:     Location{Lon: -81.431522, Lat: 44.528402},
				CreatedAt: 1591933341671,
				UpdatedAt: 1598233701672,
			},
			Speed:  25,
			Radius: 18,
		},
	}
}

func setupMinUnixNanoTime() {
	MinUnixNanoTime = time.Now().UnixNano()
}

func TestSetRadius_NormalCases(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetRadius(TestRadius)

		// assert correct location set
		assert.Equal(t, TestRadius, tc.Radius)
	}
}

func TestSetSpeed_NormalCases(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetSpeed(TestSpeed)

		// assert correct location set
		assert.Equal(t, TestSpeed, tc.Speed)
	}
}
