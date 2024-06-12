package cacheServices

import (
	"i9pxc/appTypes"
	"net/http"
	"time"
)

func staleWhileRevalidate(r *http.Request, cc *appTypes.CacheControl, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	extraTime, _ := time.ParseDuration(cc.Get("stale-while-revalidate") + "s")

	if responseIsStale(cacheData.CachedAt, cc, cacheData.Header.Get("Expires"), extraTime, false) {
		return revalidate(r, cacheData, cacheRequestKey)
	}

	return genCacheResp(cacheData)
}

func staleIfError(r *http.Request, cc *appTypes.CacheControl, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	resp, err := revalidate(r, cacheData, cacheRequestKey)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	extraTime, _ := time.ParseDuration(cc.Get("stale-if-error") + "s")

	if resp.StatusCode >= 500 && resp.StatusCode < 600 && !responseIsStale(cacheData.CachedAt, cc, cacheData.Header.Get("Expires"), extraTime, false) {

		return genCacheResp(cacheData)
	}

	return resp, err
}
