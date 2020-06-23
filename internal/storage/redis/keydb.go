package redis

import (
	"github.com/gomodule/redigo/redis"

	"location-service/internal/types"
)

type KeyDB struct {
	scripts map[string]*redis.Script
	pool    *redis.Pool
}

func NewKeyDB(pool *redis.Pool) *KeyDB {
	conn := pool.Get()
	defer conn.Close()

	ks := getKeyScripts()

	return &KeyDB{
		scripts: loadScripts(conn, ks),
		pool:    pool,
	}
}

func (db *KeyDB) Set(key, val string) error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	return err
}

func (db *KeyDB) SetIfExists(key, val string) error {
	conn := db.pool.Get()
	defer conn.Close()

	res, err := db.scripts["KEYSETEX"].Do(conn, OneKey, key, val)

	if res == nil {
		return types.ErrKeyNotFound
	}

	return err
}

func (db *KeyDB) Get(key string) (string, error) {
	conn := db.pool.Get()
	defer conn.Close()

	res, err := conn.Do("GET", key)

	if res == nil {
		return "", types.ErrKeyNotFound
	}

	return redis.String(res, err)
}

func (db *KeyDB) Delete(key string) error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (db *KeyDB) CountKeys() (int, error) {
	conn := db.pool.Get()
	defer conn.Close()

	res, err := redis.Values(conn.Do("SCAN", nil))
	if err != nil {
		return 0, err
	}

	keys, _ := redis.Strings(res[1], nil)
	return len(keys), nil
}

func (db *KeyDB) Clear() error {
	conn := db.pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")
	return err
}
