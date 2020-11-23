package tracking

import (
	"fmt"
	"testing"
	"time"

	// "time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"location-service/internal"
	"location-service/internal/mocks"
	"location-service/internal/stream"
)

var (
	service     *Service
	clientDrain *stream.Drain

	mockCourierStore *mocks.MockCourierStorer
	mockOrderStore   *mocks.MockOrderStorer

	mocker *gomock.Controller
)

func setupServiceTests() {
	mockCourierStore = mocks.NewMockCourierStorer(mocker)
	mockOrderStore = mocks.NewMockOrderStorer(mocker)

	svcDrain := stream.NewDrain()
	clientDrain = stream.NewDrain()

	svcDrain.SetInput(clientDrain.GetOutput())
	clientDrain.SetInput(svcDrain.GetOutput())

	service = NewService(mockCourierStore, mockOrderStore, svcDrain)
}

// func updateValsRandomly(tcs []internal.Courier) {
// 	for i := range tcs {
// 		tcs[i].SetLocation(rand.Float64(), rand.Float64())
// 		tcs[i].SetRadius(rand.Float64())
// 		tcs[i].SetSpeed(rand.Float64())
// 		tcs[i].SetUpdatedAt()
// 		tcs[i].SetCreatedAt()
// 	}
// }

// var (
// 	GetCourierIDNotFoundCases = []struct {
// 		inputID     string
// 		expectedErr error
// 	}{
// 		{inputID: "mockid1", expectedErr: types.ErrKeyNotFound},
// 		{inputID: "mockid2", expectedErr: types.ErrKeyNotFound},
// 		{inputID: "mockid3", expectedErr: types.ErrKeyNotFound},
// 	}
// )

// func TestGetCourierIDNotFoundCases(t *testing.T) {
// 	mocker = gomock.NewController(t)
// 	defer mocker.Finish()
// 	setupServiceTests(t)

// 	for _, c := range GetCourierIDNotFoundCases {
// 		mockCourierStore.EXPECT().GetCourier(c.inputID, &Courier{}).Return(c.expectedErr).Times(1)

// 		_, err := svc.GetCourier(c.inputID)

// 		assert.EqualError(t, err, c.expectedErr.Error())
// 	}
// }

// var (
// 	GetCourierNormalCases = []struct {
// 		inputID         string
// 		expectedCourier *Courier
// 	}{
// 		{inputID: "mockid1", expectedCourier: &Courier{}},
// 		{inputID: "mockid2", expectedCourier: &Courier{}},
// 		{inputID: "mockid3", expectedCourier: &Courier{}},
// 	}
// )

// func TestGetCourierNormalCases(t *testing.T) {
// 	mocker = gomock.NewController(t)
// 	defer mocker.Finish()
// 	setupServiceTests(t)

// 	for _, c := range GetCourierNormalCases {
// 		mockCourierStore.EXPECT().GetCourier(c.inputID, &Courier{}).Return(nil).Times(1)

// 		courier, err := svc.GetCourier(c.inputID)

// 		assert.Equal(t, courier, c.expectedCourier)
// 		assert.NoError(t, err)
// 	}
// }

// var (
// 	DeleteNormalCases = []struct {
// 		inputID     string
// 		expectedErr error
// 	}{
// 		{inputID: "mock_id1", expectedErr: types.ErrKeyNotFound},
// 		{inputID: "mock_id2", expectedErr: types.ErrKeyNotFound},
// 		{inputID: "mock_id3", expectedErr: types.ErrKeyNotFound},
// 	}

// 	DeleteIDNotFoundCases = []struct {
// 		inputID     string
// 		expectedErr error
// 	}{
// 		{inputID: "mock_id1", expectedErr: types.ErrKeyNotFound},
// 		{inputID: "mock_id2", expectedErr: types.ErrKeyNotFound},
// 		{inputID: "mock_id3", expectedErr: types.ErrKeyNotFound},
// 	}
// )

// func TestDeleteNormalCases(t *testing.T) {
// 	mocker = gomock.NewController(t)
// 	defer mocker.Finish()
// 	setupServiceTests(t)

// 	for _, c := range DeleteNormalCases {
// 		mockCourierStore.EXPECT().DeleteCourier(c.inputID).Return(nil).Times(1)

// 		err := svc.DeleteCourier(c.inputID)
// 		assert.NoError(t, err)
// 	}
// }

// func TestDeleteIDNotFoundCases(t *testing.T) {
// 	mocker = gomock.NewController(t)
// 	defer mocker.Finish()
// 	setupServiceTests(t)

// 	for _, c := range DeleteIDNotFoundCases {
// 		mockCourierStore.EXPECT().DeleteCourier(c.inputID).Return(c.expectedErr).Times(1)

// 		err := svc.DeleteCourier(c.inputID)
// 		assert.EqualError(t, err, c.expectedErr.Error())
// 	}
// }

// responseIDs: []string{"mock_order_id1"},
// 			expectedIDs: []string{"testcourierid1"},

// responseIDs: []string{"testcourierid9", "testcourierid8", "testcourierid7"},
// 			expectedIDs: []string{"testcourierid9", "testcourierid7", "testcourierid8"},

// var (
// 	GetAllNearbyOrdersNormalCases = []struct {
// 		r           *FindAllNearbyRequest
// 		t           *internal.TrackedItem
// 		expectedIDs []string
// 	}{
// 		{
// 			r: &FindAllNearbyRequest{
// 				ID:     "mock_id1",
// 				Coord:  Location{Lon: -79.481522, Lat: 43.428401},
// 				Radius: 0.0001,
// 			},
// 			t: &internal.TrackedItem{
// 				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
// 				ID:        "mock_id1",
// 				CreatedAt: 1605443925334,
// 				UpdatedAt: 1605443925334,
// 			},
// 			expectedIDs: []string{"mock_id2", "mock_id3"},
// 		},
// 		{
// 			r: &FindAllNearbyRequest{
// 				ID:     "mock_id2",
// 				Coord:  Location{Lon: -79.481322, Lat: 43.428402},
// 				Radius: 14.92,
// 			},
// 			t: &internal.TrackedItem{
// 				Coord:     internal.Location{Lon: -79.481322, Lat: 43.428402},
// 				ID:        "mock_id2",
// 				CreatedAt: 1605443925334,
// 				UpdatedAt: 1605443925334,
// 			},
// 			expectedIDs: []string{"mock_id10", "mock_id11", "mock_id12"},
// 		},
// 	}
// )

// func TestGetAllNearbyOrdersNormalCases(t *testing.T) {
// 	mocker = gomock.NewController(t)
// 	defer mocker.Finish()
// 	setupServiceTests()

// 	for _, c := range GetAllNearbyOrdersNormalCases {
// 		internal.UnixNanoNow = func(time.Time) int64 {
// 			return 1605443925334
// 		}

// 		mockOrderStore.EXPECT().FindAllNearbyOrderIDs(c.t, c.r.Radius).Return(c.expectedIDs, nil).Times(1)

// 		orderIDs, err := svc.FindAllNearbyOrderIDs(c.r)

// 		assert.NoError(t, err)
// 		assert.Equal(t, orderIDs, c.expectedIDs)
// 	}

// 	// resetting clock stub to original time function
// 	internal.UnixNanoNow = time.Time.UnixNano
// }

// var (
// 	GetAllNearbyOrdersErrorCases = []struct {
// 		inputCoord  map[string]float64
// 		inputRadius float64
// 		expectedErr error
// 	}{
// 		{
// 			inputCoord:  map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			inputRadius: 0.0001,
// 			expectedErr: types.ErrKeyNotFound,
// 		},
// 		{
// 			inputCoord:  map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			inputRadius: 14.92,
// 			expectedErr: types.ErrKeyNotFound,
// 		},
// 	}
// )

var (
	mockCoord   = internal.NewLocation(-180, 90)
	mockRadius  = 1.0001
	expectedIDs = []string{"mock_id10", "mock_id11", "mock_id12"}

	mockBadCoord  = internal.NewLocation(-180.0001, -90)
	mockBadRadius = 50.0001
)

func TestService(t *testing.T) {
	mocker = gomock.NewController(t)
	defer mocker.Finish()
	setupServiceTests()

	t.Run("when finding all nearby order ids", func(t *testing.T) {

		t.Run("should return 3 closest order ids", func(t *testing.T) {
			internal.UnixNanoNow = func(time.Time) int64 {
				return 1605443925334
			}

			mockOrderStore.
				EXPECT().
				FindAllNearbyOrderIDs(&mockCoord, mockRadius).
				Return(expectedIDs, nil).
				Times(1)

			ids, err := service.FindAllNearbyOrderIDs(&mockCoord, mockRadius)

			assert.Equal(t, expectedIDs, ids)
			assert.NoError(t, err)

			internal.UnixNanoNow = time.Time.UnixNano
		})

		t.Run("should return error when finding all nearby order ids", func(t *testing.T) {
			internal.UnixNanoNow = func(time.Time) int64 {
				return 1605443925334
			}

			mockOrderStore.
				EXPECT().
				FindAllNearbyOrderIDs(&mockCoord, mockRadius).
				Return(nil, fmt.Errorf("mock_error")).
				Times(1)

			ids, err := service.FindAllNearbyOrderIDs(&mockCoord, mockRadius)

			assert.EqualError(t, fmt.Errorf("mock_error"), err.Error())
			assert.Nil(t, ids)

			internal.UnixNanoNow = time.Time.UnixNano
		})

		t.Run("should return error when given invalid coord", func(t *testing.T) {
			internal.UnixNanoNow = func(time.Time) int64 {
				return 1605443925334
			}

			ids, err := service.FindAllNearbyOrderIDs(&mockBadCoord, mockRadius)

			assert.EqualError(t, fmt.Errorf("Lon: must be no less than -180."), err.Error())
			assert.Nil(t, ids)

			internal.UnixNanoNow = time.Time.UnixNano
		})

		// t.Run("should return error when given invalid radius", func(t *testing.T) {
		// 	internal.UnixNanoNow = func(time.Time) int64 {
		// 		return 1605443925334
		// 	}

		// 	ids, err := service.FindAllNearbyOrderIDs(&mockCoord, mockBadRadius)

		// 	assert.Equal(t, err, errors.InvalidSearchRadius)
		// 	assert.Nil(t, ids)

		// 	internal.UnixNanoNow = time.Time.UnixNano
		// })

	})

}

// var (
// 	GetAllNearbyCouriersNormalCases = []struct {
// 		idsResponse func() ([]string, error)
// 		inputCoord  map[string]float64
// 		inputradius float64
// 		expectedIDs []string
// 	}{
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid9"}, nil
// 			},
// 			map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			0.0001,
// 			[]string{"testcourierid9"},
// 		},
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid9", "testcourierid7", "testcourierid8"}, nil
// 			},
// 			map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			14.92,
// 			[]string{"testcourierid9", "testcourierid7", "testcourierid8"},
// 		},
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid9"}, nil
// 			},
// 			map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			10,
// 			[]string{"testcourierid9"},
// 		},
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid7", "testcourierid8", "testcourierid9", "testcourierid10"}, nil
// 			},
// 			map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			100,
// 			[]string{"testcourierid7", "testcourierid8", "testcourierid9", "testcourierid10"},
// 		},
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid7", "testcourierid8", "testcourierid9", "testcourierid11", "testcourierid10"}, nil
// 			},
// 			map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			198.30465,
// 			[]string{"testcourierid7", "testcourierid8", "testcourierid9", "testcourierid11", "testcourierid10"},
// 		},
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid1", "testcourierid10", "testcourierid11", "testcourierid2", "testcourierid7", "testcourierid7.5", "testcourierid8", "testcourierid9"}, nil
// 			},
// 			map[string]float64{"lon": -79.481522, "lat": 43.428401},
// 			10000,
// 			[]string{"testcourierid1", "testcourierid10", "testcourierid11", "testcourierid2", "testcourierid7", "testcourierid7.5", "testcourierid8", "testcourierid9"},
// 		},
// 		{
// 			func() ([]string, error) {
// 				return []string{"testcourierid1", "testcourierid10", "testcourierid11", "testcourierid2", "testcourierid3", "testcourierid7", "testcourierid7.5", "testcourierid8", "testcourierid9"}, nil
// 			},
// 			map[string]float64{"lon": -120.213, "lat": 0.998401},
// 			100000,
// 			[]string{"testcourierid1", "testcourierid10", "testcourierid11", "testcourierid2", "testcourierid3", "testcourierid7", "testcourierid7.5", "testcourierid8", "testcourierid9"},
// 		},
// 	}
// )

// func TestGetAllNearbyCouriersNormalCases(t *testing.T) {
// 	setupServiceTests(t)

// 	for , c := range GetAllNearbyCouriersNormalCases {
// 		// setup mock response from test dependency
// 		courierStore.GetAllNearbyFn = c.idsResponse

// 		cids, err := svc.GetAllNearbyCouriers(c.inputCoord, c.inputradius)
// 		assert.NoError(t, err)

// 		sort.Strings(c.expectedIDs)
// 		sort.Strings(cids)

// 		assert.Equal(t, c.expectedIDs, cids)
// 	}
// }

// var (
// 	UpsertCourierBadCoordCases = []struct {
// 		keyNotFoundResponse func() error
// 		badCoordResponse    func() error
// 		inputCourier        Courier
// 		expectedBadCoordErr error
// 	}{
// 		{
// 			func() error { return types.ErrKeyNotFound },
// 			func() error { return errors.Errorf("bad coord!") },
// 			Courier{
// 				TrackedItem: TrackedItem{
// 					Coord:     Location{Lon: -180.0001, Lat: 44.528402},
// 					ID:        "testcourierid1",
// 					CreatedAt: 10,
// 					UpdatedAt: 10,
// 				},
// 				Speed:  10,
// 				Radius: 10,
// 			},
// 			errors.Errorf("bad coord!"),
// 		},
// 		{
// 			func() error { return types.ErrKeyNotFound },
// 			func() error { return errors.Errorf("bad coord!") },
// 			Courier{
// 				TrackedItem: TrackedItem{
// 					Coord:     Location{Lon: -90.82, Lat: -85.05112879},
// 					ID:        "testcourierid7.5",
// 					CreatedAt: 10,
// 					UpdatedAt: 10,
// 				},
// 				Speed:  10,
// 				Radius: 10,
// 			},
// 			errors.Errorf("bad coord!"),
// 		},
// 		{
// 			func() error { return types.ErrKeyNotFound },
// 			func() error { return errors.Errorf("bad coord!") },
// 			Courier{
// 				TrackedItem: TrackedItem{
// 					Coord:     Location{Lon: 1000, Lat: 85.05112879},
// 					ID:        "testcourierid11",
// 					CreatedAt: 10,
// 					UpdatedAt: 10,
// 				},
// 				Speed:  10,
// 				Radius: 10,
// 			},
// 			errors.Errorf("bad coord!"),
// 		},
// 	}
// )

// // Tests updating courier with success in upsert (no error)
// func TestUpsertCourierUpdateCases(t *testing.T) {
// 	setupServiceTests(t)
// 	tcs := getTestCouriers()

// 	for , tc := range tcs {
// 		courierStore.UpdateFn = func() error { return nil }
// 		courierStore.GetFn = func() error { return nil }

// 		err := svc.upsertCourier(&tc)
// 		assert.NoError(t, err)

// 		, err = svc.GetCourier(tc.ID)
// 		assert.NoError(t, err)
// 	}
// }

// // Tests adding new courier with success in upsert (no error)
// func TestUpsertCourierAddCases(t *testing.T) {
// 	setupServiceTests(t)
// 	tcs := getTestCouriers()

// 	// put everything in for the first time
// 	// check that exists
// 	for , tc := range tcs {
// 		courierStore.UpdateFn = func() error { return types.ErrKeyNotFound }
// 		courierStore.AddNewFn = func() error { return nil }
// 		courierStore.GetFn = func() error { return nil }

//
// 		err := svc.upsertCourier(&tc)
// 		assert.NoError(t, err)

// 		, err = svc.GetCourier(tc.ID)
// 		assert.NoError(t, err)
// 	}
// }

// // Tests returing an error on exit when updating existing
// // courier within upsert.
// func TestUpsertCourierUpdateBadCoordCases(t *testing.T) {
// 	setupServiceTests(t)

// 	for , c := range UpsertCourierBadCoordCases {
// 		// setup mock response from test dependency
// 		courierStore.UpdateFn = c.badCoordResponse
// 		courierStore.GetFn = c.keyNotFoundResponse

//
// 		err := svc.upsertCourier(&c.inputCourier)
// 		assert.EqualError(t, err, c.expectedBadCoordErr.Error())

// 		, err = svc.GetCourier(c.inputCourier.ID)
// 		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
// 	}
// }

// // Tests returning an error on exit when adding new courier
// // within upsert.
// func TestUpsertCourierAddBadCoordCases(t *testing.T) {
// 	setupServiceTests(t)

// 	for , c := range UpsertCourierBadCoordCases {
// 		// setup mock response from test dependency
// 		courierStore.UpdateFn = c.keyNotFoundResponse
// 		courierStore.AddNewFn = c.badCoordResponse
// 		courierStore.GetFn = c.keyNotFoundResponse

//
// 		err := svc.upsertCourier(&c.inputCourier)
// 		assert.EqualError(t, err, c.expectedBadCoordErr.Error()) // change this to assert proper error!

// 		// assert key not found since nothing was inserted
// 		, err = svc.GetCourier(c.inputCourier.ID)
// 		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
// 	}
// }

// func TestTrackCourier(t *testing.T) {
// 	setupServiceTests(t)
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
