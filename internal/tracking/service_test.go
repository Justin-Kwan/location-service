package tracking

import (
	"math/rand"
	"sort"
	"testing"
	// "time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"location-service/internal/mock"
	"location-service/internal/streaming"
	"location-service/internal/types"
)

var (
	svc          *Service
	clientDrain  *streaming.Drain
	courierStore *mock.ItemStore
	orderStore   *mock.ItemStore
)

func getTestCouriers() []Courier {
	return []Courier{
		Courier{
			ID:        "test_courier_id1",
			Location:  Location{Lon: 67.1349, Lat: 78.45314},
			Speed:     45.12,
			Radius:    30,
			CreatedAt: 123456789009,
			UpdatedAt: 987654321001,
		},
		Courier{
			ID:        "test_courier_id2",
			Location:  Location{Lon: 43.9876567888854435656, Lat: 65.7654567876},
			Speed:     0.0000000000001,
			Radius:    0.0000000000002,
			CreatedAt: 987656,
			UpdatedAt: 8,
		},
		Courier{
			ID:        "test_courier_id3",
			Location:  Location{Lon: -180, Lat: -85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		Courier{
			ID:        "test_courier_id4",
			Location:  Location{Lon: 180, Lat: 85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		Courier{
			ID:        "test_courier_id5",
			Location:  Location{Lon: -180, Lat: 85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		Courier{
			ID:        "test_courier_id6",
			Location:  Location{Lon: 180, Lat: -85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		Courier{
			ID:        "test_courier_id7",
			Location:  Location{Lon: -79.661522, Lat: 43.458401},
			Speed:     50,
			Radius:    15,
			CreatedAt: 1591933701672,
			UpdatedAt: 1591933701672,
		},
		Courier{
			ID:        "test_courier_id7.5",
			Location:  Location{Lon: 0, Lat: 0},
			Speed:     0,
			Radius:    0,
			CreatedAt: 0,
			UpdatedAt: 0,
		},
		Courier{
			ID:        "test_courier_id8",
			Location:  Location{Lon: -79.661522, Lat: 43.458401},
			Speed:     50,
			Radius:    15,
			CreatedAt: 1991933701672,
			UpdatedAt: 1591933709672,
		},
		Courier{
			ID:        "test_courier_id9",
			Location:  Location{Lon: -79.481522, Lat: 43.428401},
			Speed:     50,
			Radius:    15,
			CreatedAt: 1591963701672,
			UpdatedAt: 1592933701672,
		},
		Courier{
			ID:        "test_courier_id10",
			Location:  Location{Lon: -80.481522, Lat: 43.328401},
			Speed:     20,
			Radius:    10,
			CreatedAt: 1591933701671,
			UpdatedAt: 1591933701672,
		},
		Courier{
			ID:        "test_courier_id11",
			Location:  Location{Lon: -81.431522, Lat: 44.528402},
			Speed:     25,
			Radius:    18,
			CreatedAt: 1591933341671,
			UpdatedAt: 1598233701672,
		},
	}
}

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
	for i, _ := range tcs {
		tcs[i].SetLocation(rand.Float64(), rand.Float64())
		tcs[i].SetRadius(rand.Float64())
		tcs[i].SetSpeed(rand.Float64())
		tcs[i].SetUpdatedAt()
		tcs[i].SetCreatedAt()
	}
}

var (
	GetIDNotFoundCases = []struct {
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

func TestGetIDNotFoundCases(t *testing.T) {
	setupServiceTests()

	for _, c := range GetIDNotFoundCases {
		// setup mock response from test dependency
		courierStore.GetFn = c.keyNotFoundResponse

		// function under test
		_, err := svc.GetCourier(c.inputID)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	DeleteNormalCases = []struct {
		noErrResponse       func() error
		keyNotFoundResponse func() error
		inputID             string
		expectedErr         error
	}{
		{func() error { return nil }, func() error { return types.ErrKeyNotFound }, "test_id1", types.ErrKeyNotFound},
		{func() error { return nil }, func() error { return types.ErrKeyNotFound }, "test_id2", types.ErrKeyNotFound},
		{func() error { return nil }, func() error { return types.ErrKeyNotFound }, "test_id3", types.ErrKeyNotFound},
	}

	DeleteIDNotFoundCases = []struct {
		mockResponse func() error
		inputID      string
		expectedErr  error
	}{
		{func() error { return nil }, "non_existent_id1", types.ErrKeyNotFound},
		{func() error { return nil }, "non_existent_id2", types.ErrKeyNotFound},
		{func() error { return nil }, "non_existent_id3", types.ErrKeyNotFound},
	}
)

func TestDeleteNormalCases(t *testing.T) {
	setupServiceTests()

	for _, c := range DeleteNormalCases {
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

func TestDeleteIDNotFoundCases(t *testing.T) {
	setupServiceTests()

	for _, c := range DeleteIDNotFoundCases {
		// setup mock response from test dependency
		courierStore.DeleteFn = c.mockResponse

		// function under test
		err := svc.DeleteCourier(c.inputID)
		assert.NoError(t, err)
	}
}

var (
	GetAllNearbyCouriersNormalCases = []struct {
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

func TestGetAllNearbyCouriersNormalCases(t *testing.T) {
	setupServiceTests()

	for _, c := range GetAllNearbyCouriersNormalCases {
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
	UpsertCourierBadCoordCases = []struct {
		keyNotFoundResponse func() error
		badCoordResponse    func() error
		inputCourier        Courier
		expectedBadCoordErr error
	}{
		{
			func() error { return types.ErrKeyNotFound },
			func() error { return errors.Errorf("bad coord!") },
			Courier{
				ID:        "test_courier_id1",
				Location:  Location{Lon: -180.0001, Lat: 44.528402},
				Speed:     10,
				Radius:    10,
				CreatedAt: 10,
				UpdatedAt: 10,
			},
			errors.Errorf("bad coord!"),
		},
		{
			func() error { return types.ErrKeyNotFound },
			func() error { return errors.Errorf("bad coord!") },
			Courier{
				ID:        "test_courier_id7.5",
				Location:  Location{Lon: -90.82, Lat: -85.05112879},
				Speed:     10,
				Radius:    10,
				CreatedAt: 10,
				UpdatedAt: 10,
			},
			errors.Errorf("bad coord!"),
		},
		{
			func() error { return types.ErrKeyNotFound },
			func() error { return errors.Errorf("bad coord!") },
			Courier{
				ID:        "test_courier_id11",
				Location:  Location{Lon: 1000, Lat: 85.05112879},
				Speed:     10,
				Radius:    10,
				CreatedAt: 10,
				UpdatedAt: 10,
			},
			errors.Errorf("bad coord!"),
		},
	}
)

// Tests updating courier with success in upsert (no error)
func TestUpsertCourierUpdateCases(t *testing.T) {
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

// // func (s *Service) upsertCourier(c *Courier) error {
// // 	log.Printf("UPSERTING")
// // 	err := s.courierStore.Update(c)
// //
// // 	if err != nil && err == types.ErrKeyNotFound {
// // 		return s.courierStore.AddNew(c)													// tested this exit (err)
// // 	}
// //
// // 	return err																							// tested this exit (err, normal)
// // }

// Tests returing an error on exit when updating existing
// courier within upsert.
func TestUpsertCourierUpdateBadCoordCases(t *testing.T) {
	setupServiceTests()

	for _, c := range UpsertCourierBadCoordCases {
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
func TestUpsertCourierAddBadCoordCases(t *testing.T) {
	setupServiceTests()

	for _, c := range UpsertCourierBadCoordCases {
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

// func TestTrackCourier(t *testing.T) {
// 	setupServiceTests()
//
// 	courierStore.UpdateFn = func() error { return nil }
// 	// 		courierStore.AddNewFn = func() error { return nil }
//
// 	svc.TrackCourier("id")
// 	dto := TrackCourierDTO{}
//
// 	for {
// 		time.Sleep(2 * time.Second)
// 		dto.Location.Lon = rand.Float64()
// 		dto.Location.Lat = rand.Float64()
// 		dto.Speed = rand.Float64()
// 		dto.Radius = rand.Float64()
//
// 		clientDrain.Send(dto)
// 	}
//
// }
