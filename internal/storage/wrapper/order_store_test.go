package wrapper

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"location-service/internal/storage/redis"
	"location-service/internal/testutil"
)

var (
	_orderStore *OrderStore

	mockItemStore *itemStore

	orderStoreMocker *gomock.Controller
)

func setupOrderStoreTests() {
	cfg := testutil.GetConfig()

	mockItemStore = NewItemStore(mockKeyDB, mockGeoDB, &cfg.Stores.Courier)

	_orderStore = NewOrderStore(mockItemStore)
}

func TestOrderStore(t *testing.T) {
	orderStoreMocker = gomock.NewController(t)
	defer orderStoreMocker.Finish()
	setupOrderStoreTests()

	t.Run("when finding all nearby order ids", func(t *testing.T) {

		t.Run("should return 3 closest order ids", func(t *testing.T) {
			mockGeoDB.
				EXPECT().
				BatchGetAllInRadius([]*redis.GeoQuery{
					mockGeoQueryUnmatched1,
					mockGeoQueryMatched1,
				}).
				Return(expectedIDs1, nil).
				Times(1)

			ids, err := _orderStore.FindAllNearbyOrderIDs(&mockCoord1, mockRadius1)

			assert.Equal(t, ids, expectedIDs1)
			assert.NoError(t, err)
		})

		t.Run("should return error when finding all order ids nearby invalid radius", func(t *testing.T) {
			mockGeoDB.
				EXPECT().
				BatchGetAllInRadius([]*redis.GeoQuery{
					mockGeoQueryUnmatched2,
					mockGeoQueryMatched2,
				}).
				Return(nil, fmt.Errorf("invalid radius")).
				Times(1)

			ids, err := _orderStore.FindAllNearbyOrderIDs(&mockCoord2, mockInvalidRadius1)

			assert.Nil(t, ids)
			assert.Equal(t, fmt.Errorf("invalid radius"), err)
		})

	})

}
