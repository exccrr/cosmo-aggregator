package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func Set(key string, value string, ttl time.Duration) error {
	return rdb.Set(ctx, key, value, ttl).Err()
}
