package copydb

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var removeItemScript = redis.NewScript(`
local v = redis.call('HGET', KEYS[1], '` + keyVer + `')
if redis.call('DEL', KEYS[1]) == 0 then
	return 0
end
return redis.call('HINCRBY', KEYS[1], '` + keyVer + `', v+1)`)

var updateItemScript = redis.NewScript(`
local proto_set_len = 2*ARGV[1] -- key/value pairs
local proto_unset_len = ARGV[2]
local proto_argv_offset = 2

if proto_set_len+proto_unset_len ~= table.getn(ARGV)-proto_argv_offset then
	return redis.error_reply("wrong arguments count")
end

for i=1,proto_set_len,2 do
	redis.call('HSET', KEYS[1], ARGV[proto_argv_offset+i], ARGV[proto_argv_offset+i+1])
end

proto_argv_offset = proto_argv_offset + proto_set_len

for i=1,proto_unset_len do
	redis.call('HDEL', KEYS[1], ARGV[proto_argv_offset+i])
end

return redis.call('HINCRBY', KEYS[1], '` + keyVer + `', 1)`)

func loadScripts(r Redis, s ...*redis.Script) error {
	for _, ss := range s {
		if err := ss.Load(r).Err(); err != nil {
			return errors.Wrap(err, "setup script failed")
		}
	}
	return nil
}

func setupScripts(r Redis) error {
	return loadScripts(r, removeItemScript, updateItemScript)
}
