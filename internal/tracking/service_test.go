package tracking

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"location-service/internal/mock"
	"location-service/internal/streaming"
	"location-service/internal/types"
)

var (
	svc *Service

	orderStore   *mock.ItemStore
	courierStore *mock.ItemStore

	clientDrain *streaming.Drain
)

func setupServiceTests() {
	courierStore = mock.NewItemStore()
	orderStore = mock.NewItemStore()

	svcDrain := streaming.NewDrain()
	clientDrain = streaming.NewDrain()

	svcDrain.SetInput(clientDrain.GetOutput())
	clientDrain.SetInput(svcDrain.GetOutput())

	svc = NewService(courierStore, orderStore, svcDrain)
}

func updateValsRandomly(tcs []Courier) {
	for i := range tcs {
		tcs[i].SetLocation(rand.Float64(), rand.Float64())
		tcs[i].SetRadius(rand.Float64())
		tcs[i].SetSpeed(rand.Float64())
		tcs[i].SetUpdatedAt()
		tcs[i].SetCreatedAt()
	}
}

var (
	Get_IDNotFoundCases = []struct {
		keyNotFoundResponse func() error
		inputID             string
		expectedErr         error
	}{
		{func() error { return types.ErrKeyNotFound }, "non_existent_key", types.ErrKeyNotFound},
		{func() error { return types.ErrKeyNotFound }, " ", types.ErrKeyNotFound},
		{func() error { return types.ErrKeyNotFound }, "", types.ErrKeyNotFound},
		{func() error { return types.ErrKeyNotFound }, "*", types.ErrKeyNotFound},
	}
)

func TestGet_IDNotFoundCases(t *testing.T) {
	setupServiceTests()

	for _, c := range Get_IDNotFoundCases {
		// setup mock response from test dependency
		courierStore.GetFn = c.keyNotFoundResponse

		// function under test
		_, err := svc.GetCourier(c.inputID)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	Delete_NormalCases = []struct {
		noErrResponse       func() error
		keyNotFoundResponse func() error
		inputID             string
		expectedErr         error
	}{
		{func() error { return nil }, func() error { return types.ErrKeyNotFound }, "test_id1", types.ErrKeyNotFound},
		{func() error { return nil }, func() error { return types.ErrKeyNotFound }, "test_id2", types.ErrKeyNotFound},
		{func() error { return nil }, func() error { return types.ErrKeyNotFound }, "test_id3", types.ErrKeyNotFound},
	}

	Delete_IDNotFoundCases = []struct {
		mockResponse func() error
		inputID      string
		expectedErr  error
	}{
		{func() error { return nil }, "non_existent_id1", types.ErrKeyNotFound},
		{func() error { return nil }, "non_existent_id2", types.ErrKeyNotFound},
		{func() error { return nil }, "non_existent_id3", types.ErrKeyNotFound},
	}
)

func TestDelete_NormalCases(t *testing.T) {
	setupServiceTests()

	for _, c := range Delete_NormalCases {
		// setup mock response from test dependency
		courierStore.DeleteFn = c.noErrResponse
		courierStore.GetFn = c.keyNotFoundResponse

		// function under test
		err := svc.DeleteCourier(c.inputID)
		assert.NoError(t, err)

		// assert courier is deleted
		_, err = svc.GetCourier(c.inputID)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

func TestDelete_IDNotFoundCases(t *testing.T) {
	setupServiceTests()

	for _, c := range Delete_IDNotFoundCases {
		// setup mock response from test dependency
		courierStore.DeleteFn = c.mockResponse

		// function under test
		err := svc.DeleteCourier(c.inputID)
		assert.NoError(t, err)
	}
}

var (
	GetAllNearbyCouriers_NormalCases = []struct {
		idsResponse func() ([]string, error)
		inputCoord  map[string]float64
		inputradius float64
		expectedIDs []string
	}{
		{
			func() ([]string, error) {
				return []string{"test_courier_id9"}, nil
			},
			map[string]float64{"lon": -79.481522, "lat": 43.428401},
			0.0001,
			[]string{"test_courier_id9"},
		},
		{
			func() ([]string, error) {
				return []string{"test_courier_id9", "test_courier_id7", "test_courier_id8"}, nil
			},
			map[string]float64{"lon": -79.481522, "lat": 43.428401},
			14.92,
			[]string{"test_courier_id9", "test_courier_id7", "test_courier_id8"},
		},
		{
			func() ([]string, error) {
				return []string{"test_courier_id9"}, nil
			},
			map[string]float64{"lon": -79.481522, "lat": 43.428401},
			10,
			[]string{"test_courier_id9"},
		},
		{
			func() ([]string, error) {
				return []string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id10"}, nil
			},
			map[string]float64{"lon": -79.481522, "lat": 43.428401},
			100,
			[]string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id10"},
		},
		{
			func() ([]string, error) {
				return []string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id11", "test_courier_id10"}, nil
			},
			map[string]float64{"lon": -79.481522, "lat": 43.428401},
			198.30465,
			[]string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id11", "test_courier_id10"},
		},
		{
			func() ([]string, error) {
				return []string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"}, nil
			},
			map[string]float64{"lon": -79.481522, "lat": 43.428401},
			10000,
			[]string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"},
		},
		{
			func() ([]string, error) {
				return []string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id3", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"}, nil
			},
			map[string]float64{"lon": -120.213, "lat": 0.998401},
			100000,
			[]string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id3", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"},
		},
	}
)

func TestGetAllNearbyCouriers_NormalCases(t *testing.T) {
	setupServiceTests()

	for _, c := range GetAllNearbyCouriers_NormalCases {
		// setup mock response from test dependency
		courierStore.GetAllNearbyFn = c.idsResponse

		// function under test
		cids, err := svc.GetAllNearbyCouriers(c.inputCoord, c.inputradius)
		assert.NoError(t, err)

		sort.Strings(c.expectedIDs)
		sort.Strings(cids)

		assert.Equal(t, c.expectedIDs, cids)
	}
}

var (
	UpsertCourier_BadCoordCases = []struct {
		keyNotFoundResponse func() error
		badCoordResponse    func() error
		inputCourier        Courier
		expectedBadCoordErr error
	}{
		{
			func() error { return types.ErrKeyNotFound },
			func() error { return errors.Errorf("bad coord!") },

			Courier{
				TrackedItem: TrackedItem{
					Coord:     Location{Lon: -180.0001, Lat: 44.528402},
					ID:        "test_courier_id1",
					CreatedAt: 10,
					UpdatedAt: 10,
				},
				Speed:  10,
				Radius: 10,
			},
			errors.Errorf("bad coord!"),
		},
		{
			func() error { return types.ErrKeyNotFound },
			func() error { return errors.Errorf("bad coord!") },
			Courier{
				TrackedItem: TrackedItem{
					Coord:     Location{Lon: -90.82, Lat: -85.05112879},
					ID:        "test_courier_id7.5",
					CreatedAt: 10,
					UpdatedAt: 10,
				},
				Speed:  10,
				Radius: 10,
			},
			errors.Errorf("bad coord!"),
		},
		{
			func() error { return types.ErrKeyNotFound },
			func() error { return errors.Errorf("bad coord!") },
			Courier{
				TrackedItem: TrackedItem{
					Coord:     Location{Lon: 1000, Lat: 85.05112879},
					ID:        "test_courier_id11",
					CreatedAt: 10,
					UpdatedAt: 10,
				},
				Speed:  10,
				Radius: 10,
			},
			errors.Errorf("bad coord!"),
		},
	}
)

// Tests updating courier with success in upsert (no error)
func TestUpsertCourier_UpdateCases(t *testing.T) {
	setupServiceTests()
	tcs := getTestCouriers()

	for _, tc := range tcs {
		courierStore.UpdateFn = func() error { return nil }
		courierStore.GetFn = func() error { return nil }

		// function under test
		err := svc.upsertCourier(&tc)
		assert.NoError(t, err)

		_, err = svc.GetCourier(tc.ID)
		assert.NoError(t, err)
	}
}

// Tests adding new courier with success in upsert (no error)
func TestUpsertCourierAddCases(t *testing.T) {
	setupServiceTests()
	tcs := getTestCouriers()

	// put everything in for the first time
	// check that exists
	for _, tc := range tcs {
		courierStore.UpdateFn = func() error { return types.ErrKeyNotFound }
		courierStore.AddNewFn = func() error { return nil }
		courierStore.GetFn = func() error { return nil }

		// function under test
		err := svc.upsertCourier(&tc)
		assert.NoError(t, err)

		_, err = svc.GetCourier(tc.ID)
		assert.NoError(t, err)
	}
}

// Tests returing an error on exit when updating existing
// courier within upsert.
func TestUpsertCourier_UpdateBadCoordCases(t *testing.T) {
	setupServiceTests()

	for _, c := range UpsertCourier_BadCoordCases {
		// setup mock response from test dependency
		courierStore.UpdateFn = c.badCoordResponse
		courierStore.GetFn = c.keyNotFoundResponse

		// function under test
		err := svc.upsertCourier(&c.inputCourier)
		assert.EqualError(t, err, c.expectedBadCoordErr.Error())

		_, err = svc.GetCourier(c.inputCourier.ID)
		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
	}
}

// Tests returning an error on exit when adding new courier
// within upsert.
func TestUpsertCourier_AddBadCoordCases(t *testing.T) {
	setupServiceTests()

	for _, c := range UpsertCourier_BadCoordCases {
		// setup mock response from test dependency
		courierStore.UpdateFn = c.keyNotFoundResponse
		courierStore.AddNewFn = c.badCoordResponse
		courierStore.GetFn = c.keyNotFoundResponse

		// function under test
		err := svc.upsertCourier(&c.inputCourier)
		assert.EqualError(t, err, c.expectedBadCoordErr.Error()) // change this to assert proper error!

		// assert key not found since nothing was inserted
		_, err = svc.GetCourier(c.inputCourier.ID)
		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
	}
}

func TestTrackCourier(t *testing.T) {
	setupServiceTests()

	courierStore.UpdateFn = func() error { return nil }
	// 		courierStore.AddNewFn = func() error { return nil }

	svc.TrackCourier("id")
	dto := TrackCourierDTO{}

	for {
		time.Sleep(2 * time.Second)
		dto.Location.Lon = rand.Float64()
		dto.Location.Lat = rand.Float64()
		dto.Speed = rand.Float64()
		dto.Radius = rand.Float64()

		clientDrain.Send(dto)
	}

}
