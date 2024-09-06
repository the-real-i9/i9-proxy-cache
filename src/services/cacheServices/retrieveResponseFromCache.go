package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/src/appTypes"
	"i9pxc/src/db"

	"github.com/redis/go-redis/v9"
)

func retrieveResp(cacheKey string) (appTypes.CacheData, bool) {
	dbResp, err := db.RedisDB.Get(context.Background(), cacheKey).Result()

	if err == redis.Nil {
		return appTypes.CacheData{}, false
	}

	var cacheData appTypes.CacheData

	json.Unmarshal([]byte(dbResp), &cacheData)

	return cacheData, true
}
