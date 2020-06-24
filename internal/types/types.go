package types

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	// Redis DB errors
	ErrKeyNotFound    = errors.Errorf("Key passed in does not exist")
	ErrMemberNotFound = errors.Errorf("Member passed in does not exist")
	ErrNoBatchQueries = errors.Errorf("No geo queries passed into batch operation")
)

type (
	Config struct {
		WsServer   WsServerConfig `yaml:"websocket_server"`
		RedisKeyDB RedisConfig    `yaml:"redis_keydb"`
		RedisGeoDB RedisConfig    `yaml:"redis_geodb"`
		Stores     StoresConfig   `yaml:"stores"`
	}

	WsServerConfig struct {
		ReadDeadline int    `yaml:"read_deadline"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
		MsgSizeLimit int    `yaml:"message_size_limit"`
		Addr         string `yaml:"address"`
		Path         string `yaml:"path"`
	}

	RedisConfig struct {
		IdleTimeout  int    `yaml:"idle_timeout"`
		MaxIdle      int    `yaml:"max_idle_connections"`
		MaxActive    int    `yaml:"max_active_connections"`
		Addr         string `yaml:"address"`
		Password     string `yaml:"password"`
		ConnProtocol string `yaml:"connection_protocol"`
	}

	StoresConfig struct {
		Order   StoreConfig `yaml:"order"`
		Courier StoreConfig `yaml:"courier"`
	}

	StoreConfig struct {
		MatchedKey   string `yaml:"matched_key"`
		UnmatchedKey string `yaml:"unmatched_key"`
	}
)

type (
	TrackedItem interface {
		GetID() string
		GetLon() float64
		GetLat() float64
	}

	ItemStorer interface {
		AddNew(TrackedItem) error
		Get(id string, t TrackedItem) error
		GetAllNearby(coord map[string]float64, radius float64) ([]string, error)
		GetAllNearbyUnmatched(coord map[string]float64, radius float64) ([]string, error)
		Update(t TrackedItem) error
		Delete(id string) error
		SetUnmatched(id string) error
		SetMatched(id string) error
	}

	Drain interface {
		SetInput(<-chan interface{})
		GetOutput() <-chan interface{}
		Send(interface{})
		Read() interface{}
		Close()
	}
)

func MarshalJSON(t interface{}) (string, error) {
	bytes, err := json.Marshal(t)
	return string(bytes), err
}

func UnmarshalJSON(params string, t interface{}) error {
	return json.Unmarshal([]byte(params), &t)
}
