package redis

import (
	"github.com/gomodule/redigo/redis"

	"location-service/internal/types"
)

type GeoDB struct {
	scripts map[string]*redis.Script
	pool    *redis.Pool
}

type GeoQuery struct {
	Key     string
	Member  string
	Coord   GeoPos
	Radius  float64
	FromKey string
	ToKey   string
	Unit    string
	Order   string
}

type GeoPos struct {
	Lon float64
	Lat float64
}

// Returns a new geo database access struct given a redis connection
// pool and an index name to store all longitude, latitude pairs
// within the key values store.
func NewGeoDB(pool *redis.Pool) *GeoDB {
	conn := pool.Get()
	defer conn.Close()

	gs := getGeoScripts()

	return &GeoDB{
		scripts: loadScripts(conn, gs),
		pool:    pool,
	}
}

// Set adds a new member with coordinates (longitude and latitude).
func (db *GeoDB) Set(q *GeoQuery) error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := conn.Do(
		"GEOADD",
		q.Key,
		q.Coord.Lon,
		q.Coord.Lat,
		q.Member,
	)

	return err
}

// SetIfExists sets a given member in a key with new coordinates.
func (db *GeoDB) SetIfExists(q *GeoQuery) error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := db.scripts["GEOSETEX"].Do(
		conn,
		OneKey,
		q.Key,
		q.Coord.Lon,
		q.Coord.Lat,
		q.Member,
	)

	return err
}

// BatchSetIfExists batch sets one member and with one coordinate pair
// (longitude, latitude) given any number of geo queries, only if the
// member does not already exist in the key. The operations are executed
// in a pipelined transaction.
func (db *GeoDB) BatchSetIfExists(qs ...*GeoQuery) error {
	if len(qs) == 0 {
		return types.ErrNoBatchQueries
	}

	conn := db.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")

	for _, q := range qs {
		conn.Send(
			"EVALSHA",
			db.scripts["GEOSETEX"].Hash(),
			OneKey,
			q.Key,
			q.Coord.Lon,
			q.Coord.Lat,
			q.Member,
		)
	}

	_, err := conn.Do("EXEC")
	return err
}

// Returns a member's corosponding coordinates in a map.
func (db *GeoDB) Get(q *GeoQuery) (*GeoPos, error) {
	conn := db.pool.Get()
	defer conn.Close()

	res, err := redis.Positions(conn.Do("GEOPOS", q.Key, q.Member))
	emptyCoord := res[0] == nil

	if err != nil {
		return nil, err
	}
	if emptyCoord {
		return nil, types.ErrMemberNotFound
	}

	return &GeoPos{Lon: res[0][0], Lat: res[0][1]}, nil
}

// Returns a list of all members within the given radius of given
// coordinates.
func (db *GeoDB) GetAllInRadius(q *GeoQuery) ([]string, error) {
	conn := db.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do(
		"GEORADIUS",
		q.Key,
		q.Coord.Lon,
		q.Coord.Lat,
		q.Radius,
		q.Unit,
		q.Order,
	))
}

// BatchGetAllInRadius returns a list of all members within the given
// radius of given coordinates, searching in multiple keys given the
// number of queries. The operations are executed in a non-tranactional
// pipeline.
func (db *GeoDB) BatchGetAllInRadius(qs ...*GeoQuery) ([]string, error) {
	if len(qs) == 0 {
		return nil, types.ErrNoBatchQueries
	}

	conn := db.pool.Get()
	defer conn.Close()

	for _, q := range qs {
		conn.Send(
			"GEORADIUS",
			q.Key,
			q.Coord.Lon,
			q.Coord.Lat,
			q.Radius,
			q.Unit,
			q.Order,
		)
	}

	conn.Flush()

	res, err := db.readReplies(len(qs), conn)
	return redis.Strings(res, err)
}

// readReplies returns an agreggated list interfaces returned by redis
// pipeline queries.
func (db *GeoDB) readReplies(qCount int, conn redis.Conn) ([]interface{}, error) {
	allRes := make([]interface{}, 0)

	for i := 0; i < qCount; i++ {
		res, err := conn.Receive()
		if err != nil {
			return nil, err
		}

		allRes = append(allRes, res.([]interface{})...)
	}

	return allRes, nil
}

// Moves a member from its current key into a new key. The key is
// then deleted from the old index.
func (db *GeoDB) MoveMember(q *GeoQuery) error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := db.scripts["GEOMOVE"].Do(
		conn,
		TwoKeys,
		q.FromKey,
		q.ToKey,
		q.Member,
	)

	return err
}

// Delete deletes a specific point of interest (lon, lat) by member
// in the current key.
func (db *GeoDB) Delete(q *GeoQuery) error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZREM", q.Key, q.Member)
	return err
}

// BatchDelete batch deletes one member from a one key given any number of
// geo queries. The operations are executed in a pipelined transaction.
func (db *GeoDB) BatchDelete(qs ...*GeoQuery) error {
	if len(qs) == 0 {
		return types.ErrNoBatchQueries
	}

	conn := db.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	for _, q := range qs {
		conn.Send("ZREM", q.Key, q.Member)
	}

	_, err := conn.Do("EXEC")
	return err
}

// Clears the entire key member store.
func (db *GeoDB) Clear() error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")
	return err
}
