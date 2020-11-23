package redis

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"

	"location-service/internal/errors"
)

const (
	UnitKM         = "km"
	AscendingOrder = "ASC"
)

var (
	geoDB *GeoDB

	mockConn      *redigomock.Conn
	mockGeoDbPool *redis.Pool
)

func setupGeoDBTests() func() {
	mockGeoDbPool, mockConn = getMockGeoDbPool()

	geoDB = NewGeoDB(mockGeoDbPool)
	geoDB.Clear()

	return func() {
		geoDB.Clear()
	}
}

func getMockGeoDbPool() (*redis.Pool, *redigomock.Conn) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}

	return pool, conn
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
		Key:     Key,
		Radius:  radius,
		Unit:    UnitKM,
		OrderBy: AscendingOrder,
	}
}

func NewMoveMemberQuery(member, fromKey, toKey string) *GeoQuery {
	return &GeoQuery{
		Member:  member,
		FromKey: fromKey,
		ToKey:   toKey,
	}
}

func TestGeoDB(t *testing.T) {
	setupGeoDBTests()

	t.Run("should return error if no batch queries passed in", func(t *testing.T) {
		members, err := geoDB.BatchGetAllInRadius()
		assert.Equal(t, err, errors.NoBatchQueries)
		assert.Nil(t, members)
	})

	t.Run("should return list of 3 closest order points of interest", func(t *testing.T) {
		mockConn.Clear()
		cmd := mockConn.GenericCommand("GEORADIUS").ExpectStringSlice("mock_id1", "mock_id2", "mock_id3")
		q := NewRadiusQuery("mock_key", 1.1234, -1.1234, 12)

		members, err := geoDB.BatchGetAllInRadius(q)

		assert.Equal(t, []string{"mock_id1", "mock_id2", "mock_id3"}, members)
		assert.Equal(t, 1, mockConn.Stats(cmd))
		assert.NoError(t, err)
	})

	t.Run("should return Redis error while piping query to transaction", func(t *testing.T) {
		mockConn.Clear()
		cmd := mockConn.GenericCommand("GEORADIUS").ExpectError(fmt.Errorf("mock_pipe_error"))
		q := NewRadiusQuery("mock_key", 1.1234, -1.1234, 12)

		members, err := geoDB.BatchGetAllInRadius(q)

		assert.Equal(t, fmt.Errorf("mock_pipe_error"), err)
		assert.Equal(t, 1, mockConn.Stats(cmd))
		assert.Nil(t, members)
	})

	t.Run("should return Redis error while flushing entire transaction", func(t *testing.T) {
		mockConn.Clear()
		mockConn.FlushMock = func() error {
			return fmt.Errorf("mock_flush_error")
		}

		q := NewRadiusQuery("mock_key", 1.1234, -1.1234, 12)

		members, err := geoDB.BatchGetAllInRadius(q)

		assert.Equal(t, fmt.Errorf("mock_flush_error"), err)
		assert.Nil(t, members)

		mockConn.FlushMock = nil
	})

}
