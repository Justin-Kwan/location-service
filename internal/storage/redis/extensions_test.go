package redis

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

var (
	ExpectedKeySetexScript = StripSpaces(`
		assert(#KEYS == 1, 'One key must be provided')
		local key = table.remove(KEYS, 1)
		assert(#ARGV == 1, 'One argument must be provided')
		local val = table.remove(ARGV, 1)
		local keyNotFound = redis.call('EXISTS', key) == 0
		if keyNotFound then return false end
		return redis.call('SET', key, val)
	`)

	ExpectedGeoMoveScript = StripSpaces(`
		local REDIS_MAX_LON = 180
		local REDIS_MIN_LON = -180
		local REDIS_MAX_LAT = 85.05112878
		local REDIS_MIN_LAT = -85.05112878
		assert(#KEYS == 2, 'Two keys must be provided')
		local origKey, newKey = table.remove(KEYS, 1), table.remove(KEYS, 1)
		assert(#ARGV == 1, 'One members must be provided')
		local member = table.remove(ARGV, 1)
		local coord = redis.call('GEOPOS', origKey, member)
		assert(coord[1], 'Key or member invalid')
		redis.call('ZREM', origKey, member)
		local lon = math.max(REDIS_MIN_LON, math.min(coord[1][1], REDIS_MAX_LON))
		local lat = math.max(REDIS_MIN_LAT, math.min(coord[1][2], REDIS_MAX_LAT))
		return redis.call('GEOADD', newKey, lon, lat, member)
	`)

	ExpectedGeoSetexScript = StripSpaces(`
		assert(#KEYS == 1, 'One key must be provided')
		local key = table.remove(KEYS, 1)
		assert(#ARGV == 3, 'Three arguments must be provided')
		local lon, lat, member = table.remove(ARGV, 1), table.remove(ARGV, 1),table.remove(ARGV, 1)
		local memberNotExist = redis.call('ZRANK', key, member) == false
		if memberNotExist then return false end
		return redis.pcall('GEOADD', key, lon, lat, member)
	`)
)

func StripSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func TestGeoDBExtensions(t *testing.T) {
	setupGeoDBTests()

	t.Run("should get KEYSETEX KeyDB Redis Lua script", func(t *testing.T) {
		ks := getKeyScripts()
		assert.Equal(t, ExpectedKeySetexScript, StripSpaces(ks["KEYSETEX"]))
	})

	t.Run("should get GEOMOVE GeoDB Redis Lua script", func(t *testing.T) {
		ks := getGeoScripts()
		assert.Equal(t, ExpectedGeoMoveScript, StripSpaces(ks["GEOMOVE"]))
	})

	t.Run("should get GEOSETEX GeoDB Redis Lua script", func(t *testing.T) {
		ks := getGeoScripts()
		assert.Equal(t, ExpectedGeoSetexScript, StripSpaces(ks["GEOSETEX"]))
	})

}
