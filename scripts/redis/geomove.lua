-- file: geomove.lua
-- tag:  GEOMOVE
-- moves a geo member within a sorted set to another sorted set

local REDIS_MAX_LON = 180
local REDIS_MIN_LON = -180

local REDIS_MAX_LAT = 85.05112878
local REDIS_MIN_LAT = -85.05112878

assert(#KEYS == 2, 'Two keys must be provided')
local origKey, newKey = table.remove(KEYS, 1), table.remove(KEYS, 1)

assert(#ARGV == 1, 'One members must be provided')
local member = table.remove(ARGV, 1)

-- get coordinates of existing geo member in key
local coord = redis.call('GEOPOS', origKey, member)
assert(coord[1], 'Key or member invalid')

-- delete geo member from original key
redis.call('ZREM', origKey, member)

-- filter within min/max lon and lat due to redis' float imprecision
local lon = math.max(REDIS_MIN_LON, math.min(coord[1][1], REDIS_MAX_LON))
local lat = math.max(REDIS_MIN_LAT, math.min(coord[1][2], REDIS_MAX_LAT))

-- add geo member to new key
return redis.call('GEOADD', newKey, lon, lat, member)
