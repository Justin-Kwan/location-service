package wrapper

import (
	"fmt"
	"location-service/internal"
	"location-service/internal/storage/redis"
	"location-service/internal/testutil"
	"testing"

	"location-service/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	_itemStore *itemStore

	mockKeyDB *mocks.MockKeyDB
	mockGeoDB *mocks.MockGeoDB

	mocker *gomock.Controller
)

func setupItemStoreTests() {
	cfg := testutil.GetConfig()

	mockKeyDB = mocks.NewMockKeyDB(mocker)
	mockGeoDB = mocks.NewMockGeoDB(mocker)

	_itemStore = NewItemStore(mockKeyDB, mockGeoDB, &cfg.Stores.Courier)
}

var (
	mockCoord1 = internal.NewLocation(12, 12.43)
	mockCoord2 = internal.NewLocation(-12, -12.43)

	mockRadius1        = 10.123
	mockInvalidRadius1 = -10.123

	expectedIDs1 = []string{"mock_id1", "mock_id2", "mock_id3"}
	expectedIDs2 = []string{"mock_id1", "mock_id2", "mock_id3"}

	mockGeoQueryUnmatched1 = &redis.GeoQuery{
		Coord:   redis.GeoPos{Lon: mockCoord1.Lon, Lat: mockCoord1.Lat},
		Key:     "unmatched:test_couriers",
		Radius:  mockRadius1,
		Unit:    kmUnit,
		OrderBy: orderByAsc,
	}

	mockGeoQueryMatched1 = &redis.GeoQuery{
		Coord:   redis.GeoPos{Lon: mockCoord1.Lon, Lat: mockCoord1.Lat},
		Key:     "matched:test_couriers",
		Radius:  mockRadius1,
		Unit:    kmUnit,
		OrderBy: orderByAsc,
	}

	mockGeoQueryUnmatched2 = &redis.GeoQuery{
		Coord:   redis.GeoPos{Lon: mockCoord2.Lon, Lat: mockCoord2.Lat},
		Key:     "unmatched:test_couriers",
		Radius:  mockInvalidRadius1,
		Unit:    kmUnit,
		OrderBy: orderByAsc,
	}

	mockGeoQueryMatched2 = &redis.GeoQuery{
		Coord:   redis.GeoPos{Lon: mockCoord2.Lon, Lat: mockCoord2.Lat},
		Key:     "matched:test_couriers",
		Radius:  mockInvalidRadius1,
		Unit:    kmUnit,
		OrderBy: orderByAsc,
	}
)

func TestItemStore(t *testing.T) {
	mocker = gomock.NewController(t)
	defer mocker.Finish()
	setupItemStoreTests()

	t.Run("when finding all item ids nearby given radius", func(t *testing.T) {

		t.Run("should find all item ids nearby given radius", func(t *testing.T) {
			mockGeoDB.
				EXPECT().
				BatchGetAllInRadius([]*redis.GeoQuery{
					mockGeoQueryUnmatched1,
					mockGeoQueryMatched1,
				}).
				Return(expectedIDs1, nil).
				Times(1)

			ids, err := _itemStore.findAllNearbyItemIDs(&mockCoord1, mockRadius1)

			assert.Equal(t, ids, expectedIDs1)
			assert.NoError(t, err)
		})

		t.Run("should return error when finding all item ids nearby invalid radius", func(t *testing.T) {
			mockGeoDB.
				EXPECT().
				BatchGetAllInRadius([]*redis.GeoQuery{
					mockGeoQueryUnmatched2,
					mockGeoQueryMatched2,
				}).
				Return(nil, fmt.Errorf("invalid radius")).
				Times(1)

			ids, err := _itemStore.findAllNearbyItemIDs(&mockCoord2, mockInvalidRadius1)

			assert.Nil(t, ids)
			assert.Equal(t, fmt.Errorf("invalid radius"), err)
		})

	})

}
