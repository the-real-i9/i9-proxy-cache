package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisDB *redis.Client

func InitRedisClient() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
}
