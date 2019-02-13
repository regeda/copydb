package copydb

import "github.com/go-redis/redis"

type Redis interface {
	Pipeline() redis.Pipeliner
	HGetAll(key string) *redis.StringStringMapCmd
	ZRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd
	Subscribe(channels ...string) *redis.PubSub

	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(script string) *redis.StringCmd
}
