package redis

import (
  "github.com/gomodule/redigo/redis"
)

const (
	ManualKeyCount = -1
	OneKey         = 1
	TwoKeys        = 2
)

func getKeyScripts() map[string]string {
	return map[string]string{
		"KEYSETEX": `assert(#KEYS == 1, 'One key must be provided')
  							 local key = table.remove(KEYS, 1)
  							 assert(#ARGV == 1, 'One argument must be provided')
  							 local val = table.remove(ARGV, 1)
  							 local keyNotFound = redis.call('EXISTS', key) == 0
  							 if keyNotFound then return false end
  							 return redis.call('SET', key, val)`,
	}
}

func getGeoScripts() map[string]string {
	return map[string]string{
		"GEOMOVE": `local REDIS_MAX_LON = 180
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
								 return redis.call('GEOADD', newKey, lon, lat, member)`,

		"GEOSETEX": `assert(#KEYS == 1, 'One key must be provided')
								 local key = table.remove(KEYS, 1)
								 assert(#ARGV == 3, 'Three arguments must be provided')
								 local lon, lat, member = table.remove(ARGV, 1), table.remove(ARGV, 1),table.remove(ARGV, 1)
								 local memberNotExist = redis.call('ZRANK', key, member) == false
								 if memberNotExist then return false end
								 return redis.pcall('GEOADD', key, lon, lat, member)`,
	}
}

func loadScripts(conn redis.Conn, scripts map[string]string) map[string]*redis.Script {
	cachedScripts := make(map[string]*redis.Script)

	for tag, s := range scripts {
		cachedScripts[tag] = redis.NewScript(ManualKeyCount, s)
		cachedScripts[tag].Load(conn)
	}

	return cachedScripts
}
