package redis

import (
	"log"
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

func NewPool(redisCfg *types.RedisConfig) *redis.Pool {
	cfg := setConfig(redisCfg)

	pool := &redis.Pool{
		MaxIdle:     cfg.maxIdleConn,
		MaxActive:   cfg.maxActiveConn,
		IdleTimeout: time.Duration(cfg.idleConnTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			log.Println("REDIS ADDR: ", cfg.addr)
			log.Println("REDIS PROTOCOL: ", cfg.protocol)

			conn, err := redis.Dial(cfg.protocol, cfg.addr)
			if err != nil {
				panic(err)
			}

			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")

			if err != nil {
				panic(err)
			}

			return err
		},
	}

	return pool
}

func setConfig(redisCfg *types.RedisConfig) *PoolConfig {
	return &PoolConfig{
		idleConnTimeout: redisCfg.IdleTimeout,
		maxIdleConn:     redisCfg.MaxIdle,
		maxActiveConn:   redisCfg.MaxActive,
		addr:            redisCfg.Addr,
		password:        redisCfg.Password,
		protocol:        redisCfg.Protocol,
	}
}
