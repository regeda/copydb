package copydb

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var removeItemScript = redis.NewScript(`
local v = redis.call('HGET', KEYS[1], '__ver')
if redis.call('DEL', KEYS[1]) == 0 then
	return 0
end
return redis.call('HINCRBY', KEYS[1], '__ver', v+1)`)

var updateItemScript = redis.NewScript(`
local proto_set_len = 2*ARGV[1] -- key/value pairs
local proto_unset_len = ARGV[2]
local proto_argv_offset = 2

if proto_set_len+proto_unset_len ~= table.getn(ARGV)-proto_argv_offset then
	return {err="wrong arguments count"}
end

for i=1,proto_set_len,2 do
	redis.call('HSET', KEYS[1], ARGV[proto_argv_offset+i], ARGV[proto_argv_offset+i+1])
end

proto_argv_offset = proto_argv_offset + proto_set_len

for i=1,proto_unset_len do
	redis.call('HDEL', KEYS[1], ARGV[proto_argv_offset+i])
end

return redis.call('HINCRBY', KEYS[1], '__ver', 1)`)

// KEYS:
//	- list
//	- item
// ARGV:
//	- id
//	- deadline
var purgeItemScript = redis.NewScript(`
local score = redis.call('ZSCORE', KEYS[1], ARGV[1])
if score == nil then
	return 1
end
if score < ARGV[2] then
	redis.call('DEL', KEYS[2])
	redis.call('ZREM', KEYS[1], ARGV[1])
	return 1
end
return 0`)

func setupScripts(r Redis) error {
	if err := removeItemScript.Load(r).Err(); err != nil {
		return errors.Wrap(err, "remove script")
	}
	if err := updateItemScript.Load(r).Err(); err != nil {
		return errors.Wrap(err, "update script")
	}
	if err := purgeItemScript.Load(r).Err(); err != nil {
		return errors.Wrap(err, "purge script")
	}
	return nil
}