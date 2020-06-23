-- file: keysetex.lua
-- tag:  KEYSETEX
-- only sets a key value pair if the key already exists

assert(#KEYS == 1, 'One key must be provided')
local key = table.remove(KEYS, 1)

assert(#ARGV == 1, 'One argument must be provided')
local val = table.remove(ARGV, 1)

local keyNotFound = redis.call('EXISTS', key) == 0

if keyNotFound then
  return false
end

return redis.call('SET', key, val)
