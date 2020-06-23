package wrappers

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"location-service/internal/mock"
	"location-service/internal/storage/redis"
	"location-service/internal/testutil"
	"location-service/internal/types"
)

const (
	MaxFloatDelta = 0.01001
)

var (
	itemStore *ItemStore
	keyDB     *redis.KeyDB
	geoDB     *redis.GeoDB
)

func getTestCouriers() []mock.Courier {
	return []mock.Courier{
		mock.Courier{
			ID:        "test_courier_id1",
			Location:  mock.Location{Lon: 67.1349, Lat: 78.45314},
			Speed:     45.12,
			Radius:    30,
			CreatedAt: 123456789009,
			UpdatedAt: 987654321001,
		},
		mock.Courier{
			ID:        "test_courier_id2",
			Location:  mock.Location{Lon: 43.9876567888854435656, Lat: 65.7654567876},
			Speed:     0.0000000000001,
			Radius:    0.0000000000002,
			CreatedAt: 987656,
			UpdatedAt: 8,
		},
		mock.Courier{
			ID:        "test_courier_id3",
			Location:  mock.Location{Lon: -180, Lat: -85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		mock.Courier{
			ID:        "test_courier_id4",
			Location:  mock.Location{Lon: 180, Lat: 85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		mock.Courier{
			ID:        "test_courier_id5",
			Location:  mock.Location{Lon: -180, Lat: 85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		mock.Courier{
			ID:        "test_courier_id6",
			Location:  mock.Location{Lon: 180, Lat: -85.05112878},
			Speed:     34.1231412,
			Radius:    10.14,
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		mock.Courier{
			ID:        "test_courier_id7",
			Location:  mock.Location{Lon: -79.661522, Lat: 43.458401},
			Speed:     50,
			Radius:    15,
			CreatedAt: 1591933701672,
			UpdatedAt: 1591933701672,
		},
		mock.Courier{
			ID:        "test_courier_id7.5",
			Location:  mock.Location{Lon: 0, Lat: 0},
			Speed:     0,
			Radius:    0,
			CreatedAt: 0,
			UpdatedAt: 0,
		},
		mock.Courier{
			ID:        "test_courier_id8",
			Location:  mock.Location{Lon: -79.661522, Lat: 43.458401},
			Speed:     50,
			Radius:    15,
			CreatedAt: 1991933701672,
			UpdatedAt: 1591933709672,
		},
		mock.Courier{
			ID:        "test_courier_id9",
			Location:  mock.Location{Lon: -79.481522, Lat: 43.428401},
			Speed:     50,
			Radius:    15,
			CreatedAt: 1591963701672,
			UpdatedAt: 1592933701672,
		},
		mock.Courier{
			ID:        "test_courier_id10",
			Location:  mock.Location{Lon: -80.481522, Lat: 43.328401},
			Speed:     20,
			Radius:    10,
			CreatedAt: 1591933701671,
			UpdatedAt: 1591933701672,
		},
		mock.Courier{
			ID:        "test_courier_id11",
			Location:  mock.Location{Lon: -81.431522, Lat: 44.528402},
			Speed:     25,
			Radius:    18,
			CreatedAt: 1591933341671,
			UpdatedAt: 1598233701672,
		},
	}
}

func setupCourierRepoTests() func() {
	cfg := testutil.GetConfig()
	keyDBPool := redis.NewPool(&(*cfg).RedisKeyDB)
	geoDBPool := redis.NewPool(&(*cfg).RedisGeoDB)

	keyDB = redis.NewKeyDB(keyDBPool)
	geoDB = redis.NewGeoDB(geoDBPool)

	keyDB.Clear()
	geoDB.Clear()

	itemStore = NewItemStore(keyDB, geoDB, &(*cfg).Stores.Courier)

	return func() {
		keyDB.Clear()
		geoDB.Clear()
	}
}

func populateCourierRepo(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		err := itemStore.AddNew(&tc)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// assert test courier inserted correctly
		c := &mock.Courier{}
		err = itemStore.Get(tc.GetID(), c)
		if err != nil {
			geoDB.Clear()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, tc, *c)
	}
}

func swapMatchStatusRandomly() {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		if rand.Intn(2) == 1 {
			itemStore.SetUnmatched(tc.GetID())
			continue
		}
		itemStore.SetMatched(tc.GetID())
	}
}

func updateValsRandomly(tcs []mock.Courier) {
	for i := range tcs {
		tcs[i].SetLocation(rand.Float64(), rand.Float64())
		tcs[i].SetRadius(rand.Float64())
		tcs[i].SetSpeed(rand.Float64())
		tcs[i].SetUpdatedAt()
		tcs[i].SetCreatedAt()
	}
}

func TestCourierRepoAddNewNormalCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// assert correct courier was added in key db
		cStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		c := &mock.Courier{}
		types.UnmarshalJSON(cStr, c)

		assert.Equal(t, &tc, c, "should add correct courier in key db")

		// assert correct courier was added in geo db
		coord, err := geoDB.Get(&redis.GeoQuery{Key: itemStore.config.unmatchedKey, Member: tc.GetID()})
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tc.GetLon(), coord.Lon, MaxFloatDelta)
		assert.InDelta(t, tc.GetLat(), coord.Lat, MaxFloatDelta)

		// assert courier is not in matched geodb index
		_, err = geoDB.Get(&redis.GeoQuery{Key: itemStore.config.matchedKey, Member: tc.GetID()})
		assert.EqualError(t, err, types.ErrMemberNotFound.Error())
	}
}

func TestAddAlreadyExistsCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	tcs := getTestCouriers()
	updateValsRandomly(tcs)

	for _, tc := range tcs {
		// function under test
		if err := itemStore.AddNew(&tc); err != nil {
			t.Fatalf(err.Error())
			teardownTests()
		}

		// assert correct courier was reset with updated fields in key db
		cStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		c := &mock.Courier{}
		types.UnmarshalJSON(cStr, c)

		assert.Equal(t, &tc, c)

		// assert correct courier was reset in geo db
		coord, err := geoDB.Get(&redis.GeoQuery{Key: itemStore.config.unmatchedKey, Member: tc.GetID()})
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tc.GetLon(), coord.Lon, MaxFloatDelta)
		assert.InDelta(t, tc.GetLat(), coord.Lat, MaxFloatDelta)

		// assert courier is not in matched geodb index
		_, err = geoDB.Get(&redis.GeoQuery{Key: itemStore.config.matchedKey, Member: tc.GetID()})
		assert.EqualError(t, err, types.ErrMemberNotFound.Error())
	}
}

var (
	GetidNotFoundCases = []struct {
		inputid     string
		expectedErr error
	}{
		{"non_existent_key", types.ErrKeyNotFound},
		{" ", types.ErrKeyNotFound},
		{"", types.ErrKeyNotFound},
		{"*", types.ErrKeyNotFound},
	}
)

func TestCrGetNormalCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		cStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		c := &mock.Courier{}
		types.UnmarshalJSON(cStr, c)

		assert.Equal(t, tc, *c)
	}
}

func TestGetidNotFoundCases(t *testing.T) {
	for _, c := range GetidNotFoundCases {
		// function under test
		_, err := keyDB.Get(c.inputid)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	GetAllNearbyNormalCases = []struct {
		inputCoord  map[string]float64
		inputRadius float64
		expectedIDs []string
	}{
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			0.0001,
			[]string{"test_courier_id9"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			14.92,
			[]string{"test_courier_id9", "test_courier_id7", "test_courier_id8"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			10,
			[]string{"test_courier_id9"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			100,
			[]string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id10"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			198.30465,
			[]string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id11", "test_courier_id10"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			10000,
			[]string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"}},
		{map[string]float64{"lon": -120.213, "lat": 0.998401},
			100000,
			[]string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id3", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"}},
	}

	GetUnmatchedNearbyNormalCases = []struct {
		inputCoord  map[string]float64
		inputRadius float64
		expectedIDs []string
	}{
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			0.0001,
			[]string{"test_courier_id9"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			14.92,
			[]string{"test_courier_id9", "test_courier_id7", "test_courier_id8"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			10,
			[]string{"test_courier_id9"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			100,
			[]string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id10"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			198.30465,
			[]string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id11", "test_courier_id10"}},
		{map[string]float64{"lon": -79.481522, "lat": 43.428401},
			10000,
			[]string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"}},
		{map[string]float64{"lon": -120.213, "lat": 0.998401},
			100000,
			[]string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id3", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"}},
	}
)

func TestGetAllNearbyNormalCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)

	for _, c := range GetAllNearbyNormalCases {
		// function under test
		cIDs, err := itemStore.GetAllNearby(map[string]float64{
			"lon": c.inputCoord["lon"],
			"lat": c.inputCoord["lat"],
		}, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedIDs)
		sort.Strings(cIDs)

		assert.Equal(t, c.expectedIDs, cIDs)
	}
}

func TestGetAllNearbyInBothKeysCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	swapMatchStatusRandomly()

	for _, c := range GetAllNearbyNormalCases {
		// function under test
		cIDs, err := itemStore.GetAllNearby(map[string]float64{
			"lon": c.inputCoord["lon"],
			"lat": c.inputCoord["lat"],
		}, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedIDs)
		sort.Strings(cIDs)

		assert.Equal(t, c.expectedIDs, cIDs)
	}
}

func TestGetUnmatchedNearbyNormalCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)

	for _, c := range GetUnmatchedNearbyNormalCases {
		cIDs, err := itemStore.GetUnmatchedNearby(map[string]float64{
			"lon": c.inputCoord["lon"],
			"lat": c.inputCoord["lat"],
		}, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedIDs)
		sort.Strings(cIDs)

		assert.Equal(t, c.expectedIDs, cIDs)
	}

}

func TestGetUnmatchedNearbyEmptyCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	tcs := getTestCouriers()

	// setup
	for _, c := range tcs {
		err := itemStore.SetMatched(c.ID)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}
	}

	for _, c := range GetUnmatchedNearbyNormalCases {
		cIDs, err := itemStore.GetUnmatchedNearby(map[string]float64{
			"lon": c.inputCoord["lon"],
			"lat": c.inputCoord["lat"],
		}, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, []string{}, cIDs)
	}
}

func TestUpdateNormalCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	tcs := getTestCouriers()
	updateValsRandomly(tcs)

	for _, tc := range tcs {
		// function under test
		if err := itemStore.Update(&tc); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert courier in key db
		cStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		tcStr, _ := types.MarshalJSON(tc)
		assert.Equal(t, tcStr, cStr)

		// assert courier in geo db
		coord, err := geoDB.Get(&redis.GeoQuery{Key: itemStore.config.unmatchedKey, Member: tc.GetID()})
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tc.GetLon(), coord.Lon, MaxFloatDelta)
		assert.InDelta(t, tc.GetLat(), coord.Lat, MaxFloatDelta)
	}
}

func TestUpdateKeyNotFoundCases(t *testing.T) {
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		err := itemStore.Update(&tc)
		// assert error was returned since test couriers were never added
		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
	}
}

var (
	DeleteCourierNotFoundCases = []struct {
		inputid string
	}{
		{"non_existent_member"},
		{" "},
		{""},
		{"*"},
	}
)

func TestCRDeleteNormalCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	populateCourierRepo(t)
	tcs := getTestCouriers()

	for _, tc := range tcs {
		// function under test
		if err := itemStore.Delete(tc.GetID()); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert deleted courier no longer exists
		_, err := geoDB.Get(&redis.GeoQuery{Key: itemStore.config.matchedKey, Member: tc.GetID()})
		assert.EqualError(t, err, types.ErrMemberNotFound.Error())

		_, err = geoDB.Get(&redis.GeoQuery{Key: itemStore.config.unmatchedKey, Member: tc.GetID()})
		assert.EqualError(t, err, types.ErrMemberNotFound.Error())

		_, err = keyDB.Get(tc.GetID())
		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
	}
}

func TestDeleteMemberNotFoundCases(t *testing.T) {
	teardownTests := setupCourierRepoTests()
	defer teardownTests()

	for _, c := range DeleteCourierNotFoundCases {
		// function under test
		assert.NoError(t, itemStore.Delete(c.inputid))
	}
}
