package cacheServices

import (
	"context"
	"encoding/json"
	"i9pxc/appTypes"
	"i9pxc/db"
	"io"
	"net/http"
)

func respIsCacheable(resp *http.Response) bool {
	cc := appTypes.CacheCtrl{}
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

	if !(cc.Has("public") || cc.Has("max-age") || cc.Has("s-max-age") || cc.Has("no-cache") || resp.Header.Get("Expires") != "") {
		return false
	}

	return true
}

func filterHeader(resp *http.Response) {
	connHdVals := resp.Header.Values("Connection")
	for _, hdVal := range connHdVals {
		resp.Header.Del(hdVal)
	}

	resp.Header.Del("Connection")
}

func CacheResponse(originResp *http.Response, cacheRequestKey string) {

	if !respIsCacheable(originResp) {
		return
	}

	filterHeader(originResp)

	body, _ := io.ReadAll(originResp.Body)

	cacheResp, _ := json.Marshal(map[string]any{
		"header":  originResp.Header,
		"trailer": originResp.Trailer,
		"body":    body,
	})

	db.RedisDB.Set(context.Background(), cacheRequestKey, cacheResp, 0)
}
