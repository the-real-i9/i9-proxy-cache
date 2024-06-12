package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/appTypes"
	"i9pxc/db"
	"io"
	"net/http"
	"time"
)

func respIsCacheable(resp *http.Response) bool {
	cc := appTypes.CacheControl{}
	cc.Parse(resp.Header.Values("Cache-Control"))

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

func filterHeader(header http.Header) http.Header {
	header = header.Clone()

	connHdVals := header.Values("Connection")
	for _, hdVal := range connHdVals {
		header.Del(hdVal)
	}

	header.Del("Connection")

	return header
}

func CacheResponse(originResp *http.Response, cacheRequestKey string) {
	defer originResp.Body.Close()

	if !respIsCacheable(originResp) {
		return
	}

	originResp.Header = filterHeader(originResp.Header)

	body, _ := io.ReadAll(originResp.Body)

	cacheData, _ := json.Marshal(map[string]any{
		"header":   originResp.Header,
		"trailer":  originResp.Trailer,
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
