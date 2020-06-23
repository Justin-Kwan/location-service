package tracking

import (
	"testing"
	"time"

	// "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	TestLon         = 54.4124
	TestLat         = -45.41239
	TestRadius      = 15.241512
	TestSpeed       = 5.0995
)

var (
  MinUnixNanoTime int64
)

func setupMinUnixNanoTime() {
	MinUnixNanoTime = time.Now().UnixNano()
}

func TestCourierSetLocation(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetLocation(TestLon, TestLat)

		// assert correct location set
		assert.Equal(t, TestLon, tc.Location.Lon)
		assert.Equal(t, TestLat, tc.Location.Lat)
	}
}

func TestCourierSetRadius(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetRadius(TestRadius)

		// assert correct location set
		assert.Equal(t, TestRadius, tc.Radius)
	}
}

func TestCourierSetSpeed(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetSpeed(TestSpeed)

		// assert correct location set
		assert.Equal(t, TestSpeed, tc.Speed)
	}
}

func TestCourierSetCreatedAt(t *testing.T) {
  setupMinUnixNanoTime()
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetCreatedAt()

		// assert correct location set
		assert.GreaterOrEqual(t, tc.CreatedAt, MinUnixNanoTime)
	}
}

func TestCourierSetUpdatedAt(t *testing.T) {
  setupMinUnixNanoTime()
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		tc.SetUpdatedAt()

		// assert correct location set
		assert.GreaterOrEqual(t, tc.UpdatedAt, MinUnixNanoTime)
	}
}
