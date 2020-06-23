package wrappers

import (
	"location-service/internal/storage/redis"
	"location-service/internal/types"
)

const (
	DistUnit  = "km"
	DistOrder = "ASC"
)

type ItemStore struct {
	keyDB  KeyDB
	geoDB  GeoDB
	config StoreConfig
}

type StoreConfig struct {
	matchedKey   string
	unmatchedKey string
}

func NewItemStore(keyDB KeyDB, geoDB GeoDB, cfg *types.StoreConfig) *ItemStore {
	return &ItemStore{
		keyDB:  keyDB,
		geoDB:  geoDB,
		config: setConfig(cfg),
	}
}

func setConfig(cfg *types.StoreConfig) StoreConfig {
	return StoreConfig{
		matchedKey:   cfg.MatchedKey,
		unmatchedKey: cfg.UnmatchedKey,
	}
}

// tested
func (m *ItemStore) AddNew(t types.TrackedItem) error {
	tStr, err := types.MarshalJSON(t)
	if err != nil {
		return err
	}

	if err := m.keyDB.Set(t.GetID(), tStr); err != nil {
		return err
	}

	return m.geoDB.Set(
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: t.GetLon(), Lat: t.GetLat()},
			Key:    m.config.unmatchedKey,
			Member: t.GetID(),
		})
}

// tested
func (m *ItemStore) Get(id string, t types.TrackedItem) error {
	tStr, err := m.keyDB.Get(id)
	if err != nil {
		return err
	}

	return types.UnmarshalJSON(tStr, t)
}

func (m *ItemStore) GetUnmatchedNearby(coord map[string]float64, radius float64) ([]string, error) {
	return m.geoDB.BatchGetAllInRadius(
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: coord["lon"], Lat: coord["lat"]},
			Key:    m.config.unmatchedKey,
			Radius: radius,
			Unit:   DistUnit,
			Order:  DistOrder,
		})
}

// proper args??????
// tested
func (m *ItemStore) GetAllNearby(coord map[string]float64, radius float64) ([]string, error) {
	return m.geoDB.BatchGetAllInRadius(
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: coord["lon"], Lat: coord["lat"]},
			Key:    m.config.unmatchedKey,
			Radius: radius,
			Unit:   DistUnit,
			Order:  DistOrder,
		},
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: coord["lon"], Lat: coord["lat"]},
			Key:    m.config.matchedKey,
			Radius: radius,
			Unit:   DistUnit,
			Order:  DistOrder,
		})
}

// tested
func (m *ItemStore) Update(t types.TrackedItem) error {
	tStr, err := types.MarshalJSON(t)
	if err != nil {
		return err
	}

	if err := m.keyDB.SetIfExists(t.GetID(), tStr); err != nil {
		return err
	}

	return m.geoDB.BatchSetIfExists(
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: t.GetLon(), Lat: t.GetLat()},
			Key:    m.config.unmatchedKey,
			Member: t.GetID(),
		},
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: t.GetLon(), Lat: t.GetLat()},
			Key:    m.config.matchedKey,
			Member: t.GetID(),
		})
}

// tested
func (m *ItemStore) Delete(id string) error {
	if err := m.keyDB.Delete(id); err != nil {
		return err
	}

	return m.geoDB.BatchDelete(
		&redis.GeoQuery{
			Key:    m.config.unmatchedKey,
			Member: id,
		},
		&redis.GeoQuery{
			Key:    m.config.matchedKey,
			Member: id,
		})
}

func (m *ItemStore) SetMatched(id string) error {
	return m.geoDB.MoveMember(&redis.GeoQuery{
		Member:  id,
		FromKey: m.config.unmatchedKey,
		ToKey:   m.config.matchedKey,
	})
}

func (m *ItemStore) SetUnmatched(id string) error {
	return m.geoDB.MoveMember(&redis.GeoQuery{
		Member:  id,
		FromKey: m.config.matchedKey,
		ToKey:   m.config.unmatchedKey,
	})
}
