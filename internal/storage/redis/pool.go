package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"

	"location-service/internal/types"
)

type PoolConfig struct {
	addr            string
	password        string
	protocol        string
	idleConnTimeout int
	maxIdleConn     int
	maxActiveConn   int
}

func NewPool(redisCfg types.RedisConfig) *redis.Pool {
	cfg := setConfig(redisCfg)

	return &redis.Pool{
		MaxIdle:     cfg.maxIdleConn,
		MaxActive:   cfg.maxActiveConn,
		IdleTimeout: time.Duration(cfg.idleConnTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(cfg.protocol, cfg.addr)
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func setConfig(redisCfg types.RedisConfig) PoolConfig {
	return PoolConfig{
		idleConnTimeout: redisCfg.IdleTimeout,
		maxIdleConn:     redisCfg.MaxIdle,
		maxActiveConn:   redisCfg.MaxActive,
		addr:            redisCfg.Addr,
		password:        redisCfg.Password,
		protocol:        redisCfg.Protocol,
	}
}
