package initializers

import (
	"i9pxc/globals"
	"os"

	"github.com/gofiber/storage/redis/v3"
	"github.com/joho/godotenv"
)

func initRedisStore() {
	store := redis.New()

	globals.RedisStore = store
}

func Init() error {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	initRedisStore()

	return nil
}
