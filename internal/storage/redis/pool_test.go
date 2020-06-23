package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"location-service/internal/testutil"
	"location-service/internal/types"
)

var (
	_keyDBCfg *types.RedisConfig
)

func setupPoolTests() {
	cfg := testutil.GetConfig()
	_keyDBCfg = &(*cfg).RedisKeyDB
}

func TestNewPool(t *testing.T) {
	setupPoolTests()

	pool := NewPool(_keyDBCfg)
	assert.Equal(t, 500, pool.MaxIdle, "should set max idle connections")
	assert.Equal(t, 1200, pool.MaxActive, "should set max active connections")
	assert.NotNil(t, pool.IdleTimeout, "should set idle timeout")
}
