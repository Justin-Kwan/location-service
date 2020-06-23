package wrappers

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
	Set(*redis.GeoQuery) error
	SetIfExists(*redis.GeoQuery) error
	BatchSetIfExists(...*redis.GeoQuery) error
	Get(*redis.GeoQuery) (*redis.GeoPos, error)
	GetAllInRadius(*redis.GeoQuery) ([]string, error)
	BatchGetAllInRadius(...*redis.GeoQuery) ([]string, error)
	Delete(*redis.GeoQuery) error
	BatchDelete(...*redis.GeoQuery) error
	MoveMember(*redis.GeoQuery) error
	Clear() error
}
