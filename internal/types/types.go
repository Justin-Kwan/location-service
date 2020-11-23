package types

//go:generate mockgen -destination=../mocks/mock_drain.go -package=mocks . Drain

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	// Drain errors
	ErrDrainOutputNil = errors.New("Drain ouput channel is nil")

	// Redis DB errors
	ErrKeyNotFound    = errors.Errorf("Key passed in does not exist")
	ErrMemberNotFound = errors.Errorf("Member passed in does not exist")
	ErrNoBatchQueries = errors.Errorf("No geo queries passed into batch operation")
)

type (
	Config struct {
		WsServer   WsServerConfig   `yaml:"websocket_server"`
		GrpcServer GrpcServerConfig `yaml:"grpc_server"`
		RedisKeyDB RedisConfig      `yaml:"redis_keydb"`
		RedisGeoDB RedisConfig      `yaml:"redis_geodb"`
		Stores     StoresConfig     `yaml:"stores"`
	}

	WsServerConfig struct {
		ReadDeadline int    `yaml:"read_deadline"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
		MsgSizeLimit int    `yaml:"message_size_limit"`
		Addr         string `yaml:"address"`
		Path         string `yaml:"path"`
	}

	GrpcServerConfig struct {
		Port     string `yaml:"port"`
		Protocol string `yaml:"connection_protocol"`
	}

	RedisConfig struct {
		IdleTimeout int    `yaml:"idle_timeout"`
		MaxIdle     int    `yaml:"max_idle_connections"`
		MaxActive   int    `yaml:"max_active_connections"`
		Addr        string `yaml:"address"`
		Password    string `yaml:"password"`
		Protocol    string `yaml:"connection_protocol"`
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
	Drain interface {
		SetInput(<-chan interface{})
		GetOutput() <-chan interface{}
		Send(interface{}) bool
		Read() (interface{}, bool)
		Close()
	}
)

type TrackCourierDTO struct {
	Location struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lot"`
	} `json:"location"`
	Speed  float64 `json:"speed"`
	Radius float64 `json:"radius"`
}

func MarshalJSON(t interface{}) (string, error) {
	bytes, err := json.Marshal(t)
	return string(bytes), err
}

func UnmarshalJSON(params string, t interface{}) error {
	return json.Unmarshal([]byte(params), &t)
}
