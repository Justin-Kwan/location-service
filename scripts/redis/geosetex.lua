-- file: geosetex.lua
-- tag:  GEOSETEX
-- only sets a geo member in a sorted set if the member already exists

assert(#KEYS == 1, 'One key must be provided')
local key = table.remove(KEYS, 1)

assert(#ARGV == 3, 'Three arguments must be provided')
local lon, lat, member = table.remove(ARGV, 1), table.remove(ARGV, 1),table.remove(ARGV, 1)

local memberNotFound = redis.call('ZRANK', key, member) == false

if memberNotFound then
  return false
end

return redis.call('GEOADD', key, lon, lat, member)
