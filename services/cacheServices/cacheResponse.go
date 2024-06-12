package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/appTypes"
	"i9pxc/db"
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

func CacheResponse(originResp *http.Response, cacheRequestKey string, body []byte) {
	if !respIsCacheable(originResp) {
		return
	}

	originResp.Header = filterHeader(originResp.Header)

	cacheData, _ := json.Marshal(map[string]any{
		"header":   originResp.Header,
		"body":     body,
		"cachedAt": time.Now(),
	})

	db.RedisDB.Set(context.Background(), cacheRequestKey, cacheData, 0)
}

func RefreshCacheResponse(cacheRequestKey string) {

	var cacheData map[string]any

	dbRes, _ := db.RedisDB.Get(context.Background(), cacheRequestKey).Result()
	json.Unmarshal([]byte(dbRes), &cacheData)

	cacheData["cachedAt"] = time.Now()

	cacheDataJSON, _ := json.Marshal(cacheData)

	db.RedisDB.Set(context.Background(), cacheRequestKey, cacheDataJSON, 0)
}
