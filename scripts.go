package copydb

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var removeItemScript = redis.NewScript(`
local v = redis.call('HGET', KEYS[1], '__ver')
redis.call('DEL', KEYS[1])
if v == false then
	return nil
end
return redis.call('HSET', KEYS[1], '__ver', v)`)

func setupScripts(r Redis) error {
	if err := removeItemScript.Load(r).Err(); err != nil {
		return errors.Wrap(err, "remove item")
	}

	return nil
}
