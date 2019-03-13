package testutil

import (
	"flag"
	"testing"

	"github.com/go-redis/redis"
	"github.com/regeda/copydb"
)

var redisAddr = flag.String("redis-addr", "localhost:6379", "Redis connection address")

func init() {
	flag.Parse()
}

// NewRedis returns a client connected to standalone Redis.
func NewRedis(t *testing.T) copydb.Redis {
	client := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})

	if err := client.Ping().Err(); err != nil {
		t.Fatalf("failed to ping redis addr %s: %v", *redisAddr, err)
	}

	if err := client.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush redis addr %s: %v", *redisAddr, err)
	}

	return client
}
