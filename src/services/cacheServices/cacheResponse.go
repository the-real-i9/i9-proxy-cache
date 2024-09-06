package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/src/appTypes"
	"i9pxc/src/db"
	"i9pxc/src/helpers"
	"net/http"
	"time"
)

func respIsCacheable(resp *http.Response) bool {
	cc := &appTypes.CacheControl{}
	cc.Parse(resp.Header.Get("Cache-Control"))

	if resp.StatusCode != http.StatusOK {
		return false
	}

	if cc.Has("no-store") && !(cc.Has("must-understand") && resp.StatusCode == http.StatusOK) {
		return false
	}

	if cc.Has("private") {
		return false
	}

	if resp.Header.Get("Authorization") != "" && !(cc.Has("public") || cc.Has("s-max-age") || cc.Has("must-revalidate")) {
		return false
	}

	if !(cc.Has("max-age") || cc.Has("s-max-age") || cc.Has("no-cache") || resp.Header.Get("Expires") != "") {
		return false
	}

	return true
}

func CacheResponse(originResp *http.Response, cacheKey string, body []byte) {
	if !respIsCacheable(originResp) {
		return
	}

	originResp.Header = helpers.FilterHeader(originResp.Header)

	cacheData, _ := json.Marshal(map[string]any{
		"header":   originResp.Header,
		"body":     body,
		"cachedAt": time.Now(),
	})

	db.RedisDB.Set(context.Background(), cacheKey, cacheData, 0)
}

func RefreshCacheResponse(cacheKey string) {

	var cacheData map[string]any

	dbRes, _ := db.RedisDB.Get(context.Background(), cacheKey).Result()
	json.Unmarshal([]byte(dbRes), &cacheData)

	cacheData["cachedAt"] = time.Now()

	cacheDataJSON, _ := json.Marshal(cacheData)

	db.RedisDB.Set(context.Background(), cacheKey, cacheDataJSON, 0)
}
