package redis

import (
	"sort"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"location-service/internal/types"
	"location-service/internal/testutil"
)

const (
	PrimaryTestKey   = "primary_test_Key"
	SecondaryTestKey = "secondary_test_Key"

	UnitKM         = "km"
	AscendingOrder = "ASC"

	MaxFloatDelta = 0.01001
)

var (
	geoDB *GeoDB
)

type TestPOI struct {
	ID           string
	Key          string
	Coord        map[string]float64
	UpdatedCoord map[string]float64
}

func getTestPOIs() []TestPOI {
	return []TestPOI{
		TestPOI{
			ID:           "test_ID1",
			Key:          PrimaryTestKey,
			Coord:        map[string]float64{"lon": 1, "lat": 2},
			UpdatedCoord: map[string]float64{"lon": 86.695626, "lat": 50.6213241},
		},
		TestPOI{
			ID:           "test_ID2",
			Key:          SecondaryTestKey,
			Coord:        map[string]float64{"lon": 3.141878, "lat": 5.123},
			UpdatedCoord: map[string]float64{"lon": 0, "lat": 0},
		},
		// edge case (max lon and max lat)
		TestPOI{
			ID:           "test_ID3",
			Key:          PrimaryTestKey,
			Coord:        map[string]float64{"lon": 180, "lat": 85.05112878},
			UpdatedCoord: map[string]float64{"lon": 165.16455, "lat": 5.926418},
		},
		// edge case (min lon and min lat)
		TestPOI{
			ID:           "test_ID4",
			Key:          PrimaryTestKey,
			Coord:        map[string]float64{"lon": -180, "lat": -85.05112878},
			UpdatedCoord: map[string]float64{"lon": 41.750453, "lat": -15.3076583},
		},
		// edge case (max lon and min lat)
		TestPOI{
			ID:           "test_ID5",
			Key:          SecondaryTestKey,
			Coord:        map[string]float64{"lon": 180, "lat": -85.05112878},
			UpdatedCoord: map[string]float64{"lon": -13.550453, "lat": -26.1066642},
		},
		// edge case (min lon and max lat)
		TestPOI{
			ID:           "test_ID6",
			Key:          SecondaryTestKey,
			Coord:        map[string]float64{"lon": -180, "lat": 85.05112878},
			UpdatedCoord: map[string]float64{"lon": -174.91972, "lat": 46.88101},
		},
		TestPOI{
			ID:           "test_ID7",
			Key:          PrimaryTestKey,
			Coord:        map[string]float64{"lon": 0, "lat": 0},
			UpdatedCoord: map[string]float64{"lon": 149.939399, "lat": 70.9365909},
		},
		TestPOI{
			ID:           "test_ID8",
			Key:          SecondaryTestKey,
			Coord:        map[string]float64{"lon": 45.12321, "lat": 32.124},
			UpdatedCoord: map[string]float64{"lon": 129.126293, "lat": 67.8821867},
		},
		TestPOI{
			ID:           "test_ID9",
			Key:          PrimaryTestKey,
			Coord:        map[string]float64{"lon": 75.987567, "lat": 67.124122124},
			UpdatedCoord: map[string]float64{"lon": -179.686045, "lat": 48.033385},
		},
		TestPOI{
			ID:           "test_ID10",
			Key:          SecondaryTestKey,
			Coord:        map[string]float64{"lon": 75.987, "lat": 67.12412},
			UpdatedCoord: map[string]float64{"lon": 78.211564, "lat": -84.0631939},
		},
		TestPOI{
			ID:           "test_ID11",
			Key:          SecondaryTestKey,
			Coord:        map[string]float64{"lon": 45.213, "lat": 74.98723},
			UpdatedCoord: map[string]float64{"lon": 137.216683, "lat": 63.837240},
		},
		TestPOI{
			ID:           "test_ID12",
			Key:          PrimaryTestKey,
			Coord:        map[string]float64{"lon": 45.213, "lat": 74.98723},
			UpdatedCoord: map[string]float64{"lon": 32.456583, "lat": 54.473027657},
		},
	}
}

func setupGeoDBTests() func() {
	cfg := testutil.GetConfig()
	geoDBPool := NewPool(cfg.RedisGeoDB)

	geoDB = NewGeoDB(geoDBPool)
	geoDB.Clear()

	return func() {
		geoDB.Clear()
	}
}

func NewSetQuery(Key, member string, lon, lat float64) *GeoQuery {
	return &GeoQuery{
		Coord: GeoPos{
			Lon: lon,
			Lat: lat,
		},
		Key:    Key,
		Member: member,
	}
}

func NewRadiusQuery(Key string, lon, lat, radius float64) *GeoQuery {
	return &GeoQuery{
		Coord: GeoPos{
			Lon: lon,
			Lat: lat,
		},
		Key:    Key,
		Radius: radius,
		Unit:   UnitKM,
		Order:  AscendingOrder,
	}
}

func NewMoveMemberQuery(member, fromKey, toKey string) *GeoQuery {
	return &GeoQuery{
		Member:  member,
		FromKey: fromKey,
		ToKey:   toKey,
	}
}

func populateGeoDB(t *testing.T) {
	testPOIs := getTestPOIs()

	for _, tp := range testPOIs {
		geoDB.Set(NewSetQuery(tp.Key, tp.ID, tp.Coord["lon"], tp.Coord["lat"]))

		// assert member inserted correctly
		coord, err := geoDB.Get(&GeoQuery{Key: tp.Key, Member: tp.ID})

		if err != nil {
			geoDB.Clear()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tp.Coord["lon"], coord.Lon, MaxFloatDelta)
		assert.InDelta(t, tp.Coord["lat"], coord.Lat, MaxFloatDelta)
	}
}

func oppositeTestKey(k string) string {
	if k == PrimaryTestKey {
		return SecondaryTestKey
	}
	return PrimaryTestKey
}

var (
	Set_BadCoordCases = []struct {
		inputMember string
		inputCoord  map[string]float64
		expectedErr error
	}{
		// edge cases, slightly out of lon/lat bounds
		{"test_ID1",
			map[string]float64{"lon": -180.00001, "lat": -85.05112879},
			errors.Errorf("ERR invalid longitude,latitude pair -180.000010,-85.051129")},
		{"test_ID2",
			map[string]float64{"lon": -180.00001, "lat": 76},
			errors.Errorf("ERR invalid longitude,latitude pair -180.000010,76.000000")},
		{"test_ID3",
			map[string]float64{"lon": 65, "lat": -85.05112879},
			errors.Errorf("ERR invalid longitude,latitude pair 65.000000,-85.051129")},
		{"test_ID4",
			map[string]float64{"lon": 180.00001, "lat": 85.05112879},
			errors.Errorf("ERR invalid longitude,latitude pair 180.000010,85.051129")},
		{"test_ID5",
			map[string]float64{"lon": 180.00001, "lat": 76},
			errors.Errorf("ERR invalid longitude,latitude pair 180.000010,76.000000")},
		{"test_ID6",
			map[string]float64{"lon": 65, "lat": 85.05112879},
			errors.Errorf("ERR invalid longitude,latitude pair 65.000000,85.051129")},
		{"test_ID7",
			map[string]float64{"lon": -180.00001, "lat": 85.05112879},
			errors.Errorf("ERR invalid longitude,latitude pair -180.000010,85.051129")},
		{"test_ID8",
			map[string]float64{"lon": 180.00001, "lat": -85.05112879},
			errors.Errorf("ERR invalid longitude,latitude pair 180.000010,-85.051129")},
	}
)

func TestSet_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	testPOIs := getTestPOIs()

	for _, tp := range testPOIs {
		// function under test
		err := geoDB.Set(NewSetQuery(tp.Key, tp.ID, tp.Coord["lon"], tp.Coord["lat"]))
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// get and assert point of interest exists
		coord, err := geoDB.Get(&GeoQuery{Key: tp.Key, Member: tp.ID})
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tp.Coord["lon"], coord.Lon, MaxFloatDelta)
		assert.InDelta(t, tp.Coord["lat"], coord.Lat, MaxFloatDelta)
	}
}

func TestSet_BadCoords(t *testing.T) {
	for _, c := range Set_BadCoordCases {
		err := geoDB.Set(
			NewSetQuery(PrimaryTestKey, c.inputMember, c.inputCoord["lon"], c.inputCoord["lat"]))
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	SetIfExists_MemberNotFoundCases = []struct {
		inputMember string
		inputCoord  map[string]float64
		expectedErr error
	}{
		{"non_existent_member",
			map[string]float64{"lon": 86.695626, "lat": 50.13241},
			types.ErrMemberNotFound},
		{" ",
			map[string]float64{"lon": -86765.695626, "lat": 78.6213241},
			types.ErrMemberNotFound},
		{"",
			map[string]float64{"lon": 65.695626, "lat": 23.62408},
			types.ErrMemberNotFound},
		{"*",
			map[string]float64{"lon": 87.695626, "lat": 36.4098708},
			types.ErrMemberNotFound},
	}
)

func TestGeoDBSetIfExists_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)
	testPOIs := getTestPOIs()

	for _, tp := range testPOIs {
		// function under test
		err := geoDB.SetIfExists(NewSetQuery(
			tp.Key,
			tp.ID,
			tp.UpdatedCoord["lon"],
			tp.UpdatedCoord["lat"],
		))

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert POI's Coordinates set to new values
		coord, err := geoDB.Get(&GeoQuery{Key: tp.Key, Member: tp.ID})
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tp.UpdatedCoord["lon"], coord.Lon, MaxFloatDelta)
		assert.InDelta(t, tp.UpdatedCoord["lat"], coord.Lat, MaxFloatDelta)
	}
}

func TestGeoDBSetIfExists_MemberNotFoundCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	for _, c := range SetIfExists_MemberNotFoundCases {
		// function under test
		err := geoDB.SetIfExists(NewSetQuery(
			PrimaryTestKey,
			c.inputMember,
			c.inputCoord["lon"],
			c.inputCoord["lat"],
		))

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert that no Key was added or updated since it dIDn't already exist
		_, err = geoDB.Get(&GeoQuery{Key: PrimaryTestKey, Member: c.inputMember})
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	BatchSetIfExists_NormalCases = []struct {
		inputQueries []*GeoQuery
	}{
		// batch of update queries
		{[]*GeoQuery{
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID1", Coord: GeoPos{Lon: 67.12421, Lat: 45.87676567}},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID4", Coord: GeoPos{Lon: 0, Lat: 0}},
			&GeoQuery{Key: SecondaryTestKey, Member: "test_ID10", Coord: GeoPos{Lon: 56.51243, Lat: 6.432524}},
		}},
		// batch of update queries
		{[]*GeoQuery{
			&GeoQuery{Key: SecondaryTestKey, Member: "test_ID2", Coord: GeoPos{Lon: 14, Lat: -75.43}},
			&GeoQuery{Key: SecondaryTestKey, Member: "test_ID8", Coord: GeoPos{Lon: -100.23, Lat: 5}},
			&GeoQuery{Key: SecondaryTestKey, Member: "test_ID11", Coord: GeoPos{Lon: 23, Lat: -0.004}},
			&GeoQuery{Key: SecondaryTestKey, Member: "test_ID5", Coord: GeoPos{Lon: -1.07, Lat: 3}},
			&GeoQuery{Key: SecondaryTestKey, Member: "test_ID6", Coord: GeoPos{Lon: 1, Lat: -1}},
		}},
		// batch of update queries
		{[]*GeoQuery{
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID12", Coord: GeoPos{Lon: -97.12421, Lat: 0.0089876}},
		}},
	}

	BatchSetIfExists_MemberNotFoundCases = []struct {
		inputMember string
		expectedErr error
	}{
		{"non_existent_member", types.ErrMemberNotFound},
		{" ", types.ErrMemberNotFound},
		{"", types.ErrMemberNotFound},
		{"*", types.ErrMemberNotFound},
	}
)

func TestBatchSetIfExists_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)

	// test
	for _, c := range BatchSetIfExists_NormalCases {
		err := geoDB.BatchSetIfExists(c.inputQueries...)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert multiple members are batch updated
		for _, q := range c.inputQueries {
			coord, err := geoDB.Get(&GeoQuery{Key: q.Key, Member: q.Member})
			if err != nil {
				teardownTests()
				t.Fatalf(err.Error())
			}

			assert.InDelta(t, q.Coord.Lon, coord.Lon, MaxFloatDelta)
			assert.InDelta(t, q.Coord.Lat, coord.Lat, MaxFloatDelta)
		}
	}
}

func TestBatchSetIfExists_EmptyCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	// test
	for _, c := range BatchSetIfExists_NormalCases {
		// function under test
		err := geoDB.BatchSetIfExists(c.inputQueries...)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}
	}

	// assert
	for _, c := range BatchSetIfExists_NormalCases {
		for _, q := range c.inputQueries {
			_, err := geoDB.Get(&GeoQuery{Key: q.Key, Member: q.Member})
			assert.EqualError(t, err, types.ErrMemberNotFound.Error())
		}
	}
}

func TestGeoDBGet_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	for _, c := range BatchSetIfExists_MemberNotFoundCases {
		// assert error when getting non-existent members
		q := &GeoQuery{Key: PrimaryTestKey, Member: c.inputMember}
		_, err := geoDB.Get(q)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	GetAllInRadius_NormalCases = []struct {
		inputLon        float64
		inputLat        float64
		inputRadius     float64
		expectedMembers []string
		expectedErr     error
	}{
		{90, 65, 1, []string{}, nil},
		{86.3234, 66.123, 0, []string{}, nil},
		{75.8, 67.124, 8.086, []string{"test_ID10"}, nil},
		{75.8, 67.124, 20, []string{"test_ID10", "test_ID9"}, nil},
		{4, 5, 100, []string{"test_ID2"}, nil},
		{45, 32, 19, []string{"test_ID8"}, nil},
		{46.2, 75.001, 29, []string{"test_ID11", "test_ID12"}, nil},
		{-178.8991238, -80.2312431, 536.303, []string{"test_ID4"}, nil},
		{45.12321, 32.124, 1000000, []string{"test_ID8", "test_ID10", "test_ID9", "test_ID11", "test_ID12", "test_ID2", "test_ID1", "test_ID7", "test_ID4"}, nil},
	}

	GetAllInRadius_BadCoordCases = []struct {
		inputLon        float64
		inputLat        float64
		inputRadius     float64
		expectedMembers []string
		expectedErr     error
	}{
		{90, 65, -1, []string(nil), errors.Errorf("ERR radius cannot be negative")},
		{181, 64, 40, []string(nil), errors.Errorf("ERR invalid longitude,latitude pair 181.000000,64.000000")},
		{30, -86, 10, []string(nil), errors.Errorf("ERR invalid longitude,latitude pair 30.000000,-86.000000")},
	}

	BatchGetAllInRadius_BadCoordCases = []struct {
		inputLon        float64
		inputLat        float64
		inputRadius     float64
		expectedMembers []string
		expectedErr     error
	}{
		{90, 65, -1, []string(nil), errors.Errorf("ERR radius cannot be negative")},
		{181, 64, 40, []string(nil), errors.Errorf("ERR invalid longitude,latitude pair 181.000000,64.000000")},
		{30, -86, 10, []string(nil), errors.Errorf("ERR invalid longitude,latitude pair 30.000000,-86.000000")},
	}
)

func TestGetAllInRadius_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	testPOIs := getTestPOIs()

	// setup by inserting all POIs all in one Key so they can all be geo indexed
	for _, tp := range testPOIs {
		err := geoDB.Set(NewSetQuery(
			PrimaryTestKey,
			tp.ID,
			tp.Coord["lon"],
			tp.Coord["lat"],
		))

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}
	}

	for _, c := range GetAllInRadius_NormalCases {
		// function under test
		members, err := geoDB.BatchGetAllInRadius(
			NewRadiusQuery(
				PrimaryTestKey,
				c.inputLon,
				c.inputLat,
				c.inputRadius,
			))

		// assert list of closest point of interests' Key IDs within radius is returned
		assert.Equal(t, c.expectedMembers, members)

		// if error should be returned
		if c.expectedErr != nil {
			assert.EqualError(t, err, c.expectedErr.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGetAllInRadius_BadCoordCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	testPOIs := getTestPOIs()

	// setup by inserting all POIs all in one Key so they can all be geo indexed
	for _, tp := range testPOIs {
		err := geoDB.Set(NewSetQuery(
			PrimaryTestKey,
			tp.ID,
			tp.Coord["lon"],
			tp.Coord["lat"],
		))

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}
	}

	for _, c := range GetAllInRadius_BadCoordCases {
		// function under test
		members, err := geoDB.GetAllInRadius(
			NewRadiusQuery(
				PrimaryTestKey,
				c.inputLon,
				c.inputLat,
				c.inputRadius,
			))

		// assert list of closest point of interests' Key IDs within radius is returned
		assert.Equal(t, c.expectedMembers, members)
		// assert error returned from bad Coords passed in
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

func GetAllInRadius_EmptyCase(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	// test on empty db
	members, err := geoDB.GetAllInRadius(NewRadiusQuery(
		PrimaryTestKey,
		53.123,
		84.9823,
		100000,
	))

	if err != nil {
		teardownTests()
		t.Fatalf(err.Error())
	}

	assert.Equal(t, members, []string{}, "should return list of closest POI IDs within radius")
}

func TestBatchGetAllInRadius_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)

	for _, c := range GetAllInRadius_NormalCases {
		// function under test
		members, err := geoDB.BatchGetAllInRadius(
			NewRadiusQuery(
				PrimaryTestKey,
				c.inputLon,
				c.inputLat,
				c.inputRadius,
			),
			NewRadiusQuery(
				SecondaryTestKey,
				c.inputLon,
				c.inputLat,
				c.inputRadius,
			))

		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		sort.Strings(c.expectedMembers)
		sort.Strings(members)

		// assert list of closest point of interests' Key IDs within radius is returned
		assert.Equal(t, c.expectedMembers, members)
	}
}

func TestBatchGetAllInRadius_BadCoordCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)

	for _, c := range BatchGetAllInRadius_BadCoordCases {
		// function under test
		members, err := geoDB.BatchGetAllInRadius(
			NewRadiusQuery(
				PrimaryTestKey,
				c.inputLon,
				c.inputLat,
				c.inputRadius,
			),
			NewRadiusQuery(
				SecondaryTestKey,
				c.inputLon,
				c.inputLat,
				c.inputRadius,
			))

		sort.Strings(c.expectedMembers)
		sort.Strings(members)

		// assert list of closest point of interests' Key IDs within radius is returned
		assert.Equal(t, c.expectedMembers, members)
		// assert error is returned from bad Coords
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

func TestDelete_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)
	testPOIs := getTestPOIs()

	for _, tp := range testPOIs {
		// function under test
		if err := geoDB.Delete(&GeoQuery{Key: tp.Key, Member: tp.ID}); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert point of interest is deleted
		_, err := geoDB.Get(&GeoQuery{Key: tp.Key, Member: tp.ID})
		assert.EqualError(t, err, types.ErrMemberNotFound.Error())
	}
}

var (
	BatchDelete_NormalCases = []struct {
		inputQueries []*GeoQuery
	}{
		{[]*GeoQuery{
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID1"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID5"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID8"},
		}},
		{[]*GeoQuery{
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID2"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID9"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID4"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID6"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID10"},
		}},
		{[]*GeoQuery{
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID3"},
		}},
		// already deleted case (should not return error)
		{[]*GeoQuery{
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID5"},
			&GeoQuery{Key: PrimaryTestKey, Member: "test_ID8"},
		}},
	}
)

func TestBatchDelete_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)

	// test
	for _, c := range BatchDelete_NormalCases {
		// function under test
		if err := geoDB.BatchDelete(c.inputQueries...); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert multiple members batch deleted are gone
		for _, q := range c.inputQueries {
			_, err := geoDB.Get(&GeoQuery{Key: PrimaryTestKey, Member: q.Member})
			assert.EqualError(t, err, types.ErrMemberNotFound.Error())
		}
	}
}

func TestBatchDelete_NoQueryCases(t *testing.T) {
	// no geo queries passed in test case (should return error)
	err := geoDB.BatchDelete([]*GeoQuery{}...)
	assert.EqualError(t, err, types.ErrNoBatchQueries.Error())
	err = geoDB.BatchDelete()
	assert.EqualError(t, err, types.ErrNoBatchQueries.Error())
}

func TestMoveMember_NormalCases(t *testing.T) {
	teardownTests := setupGeoDBTests()
	defer teardownTests()

	populateGeoDB(t)
	testPOIs := getTestPOIs()

	for _, tp := range testPOIs {
		// function under test
		err := geoDB.MoveMember(NewMoveMemberQuery(tp.ID, tp.Key, oppositeTestKey(tp.Key)))
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert in new index
		coord, err := geoDB.Get(&GeoQuery{Key: oppositeTestKey(tp.Key), Member: tp.ID})
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.InDelta(t, tp.Coord["lon"], coord.Lon, MaxFloatDelta, "asserting Key exists in new index")
		assert.InDelta(t, tp.Coord["lat"], coord.Lat, MaxFloatDelta, "asserting Key exists in new index")

		// assert point of interest is deleted from old index
		_, err = geoDB.Get(&GeoQuery{Key: tp.Key, Member: tp.ID})
		assert.EqualError(t, err, types.ErrMemberNotFound.Error())
	}
}
