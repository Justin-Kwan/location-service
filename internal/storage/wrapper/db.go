package wrapper

//go:generate mockgen -destination=../../mocks/mock_keydb.go -package=mocks . KeyDB
//go:generate mockgen -destination=../../mocks/mock_geodb.go -package=mocks . GeoDB

import (
	"location-service/internal/storage/redis"
)

type KeyDB interface {
	Set(key string, val string) error
	SetIfExists(key, val string) error
	Get(key string) (string, error)
	Delete(key string) error
	Clear() error
}

type GeoDB interface {
	Set(q *redis.GeoQuery) error
	SetIfExists(q *redis.GeoQuery) error
	BatchSetIfExists(qs ...*redis.GeoQuery) error
	Get(q *redis.GeoQuery) (*redis.GeoPos, error)
	GetAllInRadius(q *redis.GeoQuery) ([]string, error)
	BatchGetAllInRadius(qs ...*redis.GeoQuery) ([]string, error)
	Delete(q *redis.GeoQuery) error
	BatchDelete(qs ...*redis.GeoQuery) error
	MoveMember(q *redis.GeoQuery) error
	Clear() error
}
