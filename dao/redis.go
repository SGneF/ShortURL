package dao

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(addr, password string, db int) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	log.Println("redis connected")
}

func CacheGet(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

func CacheSet(ctx context.Context, key, value string, ttl time.Duration) error {
	return RDB.Set(ctx, key, value, ttl).Err()
}

const ShortURLKeyPrefix = "surl:"
