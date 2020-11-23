package integration

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"location-service/internal"
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

func getTestTrackedItems() []*internal.TrackedItem {
	return []*internal.TrackedItem{
		&internal.TrackedItem{
			ID:        "test_courier_id1",
			Coord:     internal.Location{Lon: 67.1349, Lat: 78.45314},
			CreatedAt: 123456789009,
			UpdatedAt: 987654321001,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id2",
			Coord:     internal.Location{Lon: 43.9876567888854435656, Lat: 65.7654567876},
			CreatedAt: 987656,
			UpdatedAt: 8,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id3",
			Coord:     internal.Location{Lon: -180, Lat: -85.05112878},
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id4",
			Coord:     internal.Location{Lon: 180, Lat: 85.05112878},
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id5",
			Coord:     internal.Location{Lon: -180, Lat: 85.05112878},
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id6",
			Coord:     internal.Location{Lon: 180, Lat: -85.05112878},
			CreatedAt: 987656765678,
			UpdatedAt: 8765676567876,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id7",
			Coord:     internal.Location{Lon: -79.661522, Lat: 43.458401},
			CreatedAt: 1591933701672,
			UpdatedAt: 1591933701672,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id7.5",
			Coord:     internal.Location{Lon: 0, Lat: 0},
			CreatedAt: 0,
			UpdatedAt: 0,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id8",
			Coord:     internal.Location{Lon: -79.661522, Lat: 43.458401},
			CreatedAt: 1991933701672,
			UpdatedAt: 1591933709672,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id9",
			Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
			CreatedAt: 1591963701672,
			UpdatedAt: 1592933701672,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id10",
			Coord:     internal.Location{Lon: -80.481522, Lat: 43.328401},
			CreatedAt: 1591933701671,
			UpdatedAt: 1591933701672,
		},
		&internal.TrackedItem{
			ID:        "test_courier_id11",
			Coord:     internal.Location{Lon: -81.431522, Lat: 44.528402},
			CreatedAt: 1591933341671,
			UpdatedAt: 1598233701672,
		},
	}
}

func setupItemStoreTests() func() {
	cfg := testutil.GetConfig()
	keyDBPool := redis.NewPool(&cfg.RedisKeyDB)
	geoDBPool := redis.NewPool(&cfg.RedisGeoDB)

	keyDB = redis.NewKeyDB(keyDBPool)
	geoDB = redis.NewGeoDB(geoDBPool)

	keyDB.Clear()
	geoDB.Clear()

	itemStore = NewItemStore(keyDB, geoDB, &cfg.Stores.Courier)

	return func() {
		keyDB.Clear()
		geoDB.Clear()
	}
}

func populateItemStore(t *testing.T) {
	tcs := getTestTrackedItems()

	for _, tc := range tcs {
		err := itemStore.AddNewItem(tc)
		if err != nil {
			t.Fatalf(err.Error())
		}

		ti, err := itemStore.GetItem(tc.GetID())
		if err != nil {
			geoDB.Clear()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, ti, tc)
	}
}

func swapMatchStatusRandomly() {
	tcs := getTestTrackedItems()

	for _, tc := range tcs {
		if rand.Intn(2) == 1 {
			itemStore.SetUnmatched(tc.GetID())
			continue
		}
		itemStore.SetMatched(tc.GetID())
	}
}

func updateValsRandomly(tcs []*internal.TrackedItem) {
	for i := range tcs {
		tcs[i].SetLocation(rand.Float64(), rand.Float64())
		tcs[i].SetUpdatedAt()
		tcs[i].SetCreatedAt()
	}
}

func TestAddNewNormalCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	tcs := getTestTrackedItems()

	for _, tc := range tcs {
		// assert correct courier was added in key db
		cStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		ti := &internal.TrackedItem{}
		types.UnmarshalJSON(cStr, ti)

		assert.Equal(t, tc, ti)

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

func TestAddNewAlreadyExistsCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	tcs := getTestTrackedItems()
	updateValsRandomly(tcs)

	for _, tc := range tcs {
		// function under test
		if err := itemStore.AddNewItem(tc); err != nil {
			t.Fatalf(err.Error())
			teardownTests()
		}

		// assert correct courier was reset with updated fields in key db
		cStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		ti := &internal.TrackedItem{}
		types.UnmarshalJSON(cStr, ti)

		assert.Equal(t, tc, ti)

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
	Get_IDNotFoundCases = []struct {
		inputid     string
		expectedErr error
	}{
		{"non_existent_key", types.ErrKeyNotFound},
		{" ", types.ErrKeyNotFound},
		{"", types.ErrKeyNotFound},
		{"*", types.ErrKeyNotFound},
	}
)

func TestGet_NormalCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	tcs := getTestTrackedItems()

	for _, tc := range tcs {
		// function under test
		tiStr, err := keyDB.Get(tc.GetID())
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		ti := &internal.TrackedItem{}
		types.UnmarshalJSON(tiStr, ti)

		assert.Equal(t, ti, tc)
	}
}

func TestGet_IDNotFoundCases(t *testing.T) {
	for _, c := range Get_IDNotFoundCases {
		// function under test
		_, err := keyDB.Get(c.inputid)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	FindAllNearbyIDsNormalCases = []struct {
		inputTrackedItem *internal.TrackedItem
		inputRadius      float64
		expectedIDs      []string
	}{
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id1",
				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 0.0001,
			expectedIDs: []string{"test_courier_id9"},
		},
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id2",
				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 14.92,
			expectedIDs: []string{"test_courier_id9", "test_courier_id7", "test_courier_id8"},
		},
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id3",
				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 10,
			expectedIDs: []string{"test_courier_id9"},
		},
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id4",
				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 100,
			expectedIDs: []string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id10"},
		},
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id5",
				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 198.30465,
			expectedIDs: []string{"test_courier_id7", "test_courier_id8", "test_courier_id9", "test_courier_id11", "test_courier_id10"},
		},
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id6",
				Coord:     internal.Location{Lon: -79.481522, Lat: 43.428401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 10000,
			expectedIDs: []string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"},
		},
		{
			inputTrackedItem: &internal.TrackedItem{
				ID:        "mock_id7",
				Coord:     internal.Location{Lon: -120.213, Lat: 0.998401},
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			inputRadius: 100000,
			expectedIDs: []string{"test_courier_id1", "test_courier_id10", "test_courier_id11", "test_courier_id2", "test_courier_id3", "test_courier_id7", "test_courier_id7.5", "test_courier_id8", "test_courier_id9"},
		},
	}
)

func TestFindAllNearbyIDsNormalCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()
	populateItemStore(t)

	for _, c := range FindAllNearbyIDsNormalCases {
		tiIDs, err := itemStore.findAllNearbyItemIDs(c.inputTrackedItem, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedIDs)
		sort.Strings(tiIDs)

		assert.Equal(t, tiIDs, c.expectedIDs, c.inputTrackedItem.ID)
	}
}

// Tests getting all tracked item IDs that may exist in either
// matched or unmatched Redis keyspace
func TestFindAllNearbyIDsInBothKeysCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	swapMatchStatusRandomly()

	for _, c := range FindAllNearbyIDsNormalCases {
		tiIDs, err := itemStore.findAllNearbyItemIDs(c.inputTrackedItem, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedIDs)
		sort.Strings(tiIDs)

		assert.Equal(t, tiIDs, c.expectedIDs)
	}
}

var (
	GetUnmatchedNearby_NormalCases = []struct {
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

func TestGetUnmatchedNearbyNormalCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)

	for _, c := range GetUnmatchedNearby_NormalCases {
		tiIDs, err := itemStore.GetUnmatchedNearby(map[string]float64{
			"lon": c.inputCoord["lon"],
			"lat": c.inputCoord["lat"],
		}, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedIDs)
		sort.Strings(tiIDs)

		assert.Equal(t, tiIDs, c.expectedIDs)
	}

}

func TestGetUnmatchedNearby_EmptyCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	tcs := getTestTrackedItems()

	// setup
	for _, c := range tcs {
		err := itemStore.SetMatched(c.ID)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}
	}

	for _, c := range GetUnmatchedNearby_NormalCases {
		tiIDs, err := itemStore.GetUnmatchedNearby(map[string]float64{
			"lon": c.inputCoord["lon"],
			"lat": c.inputCoord["lat"],
		}, c.inputRadius)

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, tiIDs, []string{})
	}
}

func TestUpdate_NormalCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	tcs := getTestTrackedItems()
	updateValsRandomly(tcs)

	for _, tc := range tcs {
		// function under test
		if err := itemStore.Update(tc); err != nil {
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

func TestUpdate_KeyNotFoundCases(t *testing.T) {
	tcs := getTestTrackedItems()

	for _, tc := range tcs {
		// function under test
		err := itemStore.Update(tc)
		// assert error was returned since test couriers were never added
		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
	}
}

var (
	Delete_TrackedItemNotFoundCases = []struct {
		inputid string
	}{
		{"non_existent_member"},
		{" "},
		{""},
		{"*"},
	}
)

func TestDelete_NormalCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	populateItemStore(t)
	tcs := getTestTrackedItems()

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

func TestDelete_MemberNotFoundCases(t *testing.T) {
	teardownTests := setupItemStoreTests()
	defer teardownTests()

	for _, c := range Delete_TrackedItemNotFoundCases {
		// function under test
		assert.NoError(t, itemStore.Delete(c.inputid))
	}
}
