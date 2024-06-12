package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/appTypes"
	"i9pxc/db"

	"github.com/redis/go-redis/v9"
)

func retrieveResp(cacheRequestKey string) (appTypes.CacheData, bool) {
	ctx := context.Background()

	var cacheData appTypes.CacheData

	dbResp, err := db.RedisDB.Get(ctx, cacheRequestKey).Result()

	if err == redis.Nil {
		return appTypes.CacheData{}, false
	}

	json.Unmarshal([]byte(dbResp), &cacheData)

	return cacheData, true
}
