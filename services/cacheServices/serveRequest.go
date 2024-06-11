package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/appTypes"
	"i9pxc/db"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func ServeRequest(r *http.Request, cacheRequestKey string) (appTypes.CacheResp, bool) {
	ctx := context.Background()

	dbResp, err := db.RedisDB.Get(ctx, cacheRequestKey).Result()
	if err == redis.Nil {
		return appTypes.CacheResp{}, false
	}

	var cacheResp appTypes.CacheResp

	json.Unmarshal([]byte(dbResp), &cacheResp)

	return cacheResp, true
}
