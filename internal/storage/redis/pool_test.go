package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"location-service/internal/testutil"
)

func TestNewPool(t *testing.T) {

	t.Run("should create new connection pool with correct configurations", func(t *testing.T) {
		cfg := testutil.GetConfig()

		pool := NewPool(&cfg.RedisKeyDB)
		assert.Equal(t, 500, pool.MaxIdle, "should set max idle connections")
		assert.Equal(t, 1200, pool.MaxActive, "should set max active connections")
		assert.NotNil(t, pool.IdleTimeout, "should set idle timeout")
	})

}
