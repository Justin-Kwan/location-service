package wrapper

import (
	"location-service/internal"
	"location-service/internal/storage/redis"
	"location-service/internal/types"
	"log"
)

const (
	kmUnit     = "km"
	orderByAsc = "ASC"
)

type ItemStore struct {
	config *storeConfig
	keyDB  KeyDB
	geoDB  GeoDB
}

// TODO: store unit of measurement in config
type storeConfig struct {
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

func setConfig(cfg *types.StoreConfig) *storeConfig {
	return &storeConfig{
		matchedKey:   cfg.MatchedKey,
		unmatchedKey: cfg.UnmatchedKey,
	}
}

// AddNewItem ...
func (m *ItemStore) addNewItem(t *internal.TrackedItem) error {
	tStr, err := types.MarshalJSON(t)
	if err != nil {
		return err
	}

	if err := m.keyDB.Set(t.GetID(), tStr); err != nil {
		return err
	}

	return m.geoDB.Set(
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: t.Coord.Lon, Lat: t.Coord.Lat},
			Key:    m.config.unmatchedKey,
			Member: t.GetID(),
		})
}

// GetItem ...
func (m *ItemStore) getItem(id string) (*internal.TrackedItem, error) {
	t := &internal.TrackedItem{}

	tStr, err := m.keyDB.Get(id)
	if err != nil {
		return t, err
	}

	err = types.UnmarshalJSON(tStr, t)
	return t, err
}

// GetUnmatchedNearby ...
func (m *ItemStore) getUnmatchedNearby(coord map[string]float64, radius float64) ([]string, error) {
	itemIDs, err := m.geoDB.BatchGetAllInRadius(
		&redis.GeoQuery{
			Coord:   redis.GeoPos{Lon: coord["lon"], Lat: coord["lat"]},
			Key:     m.config.unmatchedKey,
			Radius:  radius,
			Unit:    kmUnit,
			OrderBy: orderByAsc,
		})

	return itemIDs, err
}

// FindAllNearbyItemIDs returns all IDs of items within
// a given radius of a coordinate (of trackable item).
func (m *ItemStore) findAllNearbyItemIDs(coord *internal.Location, radius float64) ([]string, error) {
	log.Println("<<<<<<<<<<<<<<< ITEM STORE >>>>>>>>>>")
	log.Println("COORD LON REQUESTED", coord.Lon)
	log.Println("COORD LAT REQUESTED", coord.Lat)
	log.Println("SEARCH RADIUS", radius)

	itemIDs, err := m.geoDB.BatchGetAllInRadius(
		&redis.GeoQuery{
			Coord:   redis.GeoPos{Lon: coord.Lon, Lat: coord.Lat},
			Key:     m.config.unmatchedKey,
			Radius:  radius,
			Unit:    kmUnit,
			OrderBy: orderByAsc,
		},
		&redis.GeoQuery{
			Coord:   redis.GeoPos{Lon: coord.Lon, Lat: coord.Lat},
			Key:     m.config.matchedKey,
			Radius:  radius,
			Unit:    kmUnit,
			OrderBy: orderByAsc,
		})

	return itemIDs, err
}

// Update ...
func (m *ItemStore) update(t *internal.TrackedItem) error {
	tStr, err := types.MarshalJSON(t)
	if err != nil {
		return err
	}

	if err := m.keyDB.SetIfExists(t.ID, tStr); err != nil {
		return err
	}

	return m.geoDB.BatchSetIfExists(
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: t.Coord.Lon, Lat: t.Coord.Lat},
			Key:    m.config.unmatchedKey,
			Member: t.GetID(),
		},
		&redis.GeoQuery{
			Coord:  redis.GeoPos{Lon: t.Coord.Lon, Lat: t.Coord.Lat},
			Key:    m.config.matchedKey,
			Member: t.GetID(),
		})
}

// Delete ...
func (m *ItemStore) delete(id string) error {
	if err := m.keyDB.Delete(id); err != nil {
		return err
	}

	return m.geoDB.BatchDelete(
		&redis.GeoQuery{Key: m.config.unmatchedKey, Member: id},
		&redis.GeoQuery{Key: m.config.matchedKey, Member: id})
}

// SetMatched ...
func (m *ItemStore) setMatched(id string) error {
	return m.geoDB.MoveMember(&redis.GeoQuery{
		Member:  id,
		FromKey: m.config.unmatchedKey,
		ToKey:   m.config.matchedKey,
	})
}

// SetUnmatched ...
func (m *ItemStore) setUnmatched(id string) error {
	return m.geoDB.MoveMember(&redis.GeoQuery{
		Member:  id,
		FromKey: m.config.matchedKey,
		ToKey:   m.config.unmatchedKey,
	})
}
