package cacheServices

import (
	"fmt"
	"i9pxc/appTypes"
	"net/http"
)

func staleWhileRevalidate(r *http.Request, cc appTypes.CacheControl, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	var revGraceSec float64

	fmt.Sscanf(cc.Get("stale-while-revalidate"), "%f", &revGraceSec)

	if responseIsStale(cacheData.CachedAt, getMaxAge(cc)+revGraceSec) {
		return revalidate(r, cacheData, cacheRequestKey)
	}

	return genCacheResp(cacheData.Body)
}

func staleIfError(r *http.Request, cc appTypes.CacheControl, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	resp, err := revalidate(r, cacheData, cacheRequestKey)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	var errGraceSec float64

	fmt.Sscanf(cc.Get("stale-if-error"), "%f", &errGraceSec)

	if resp.StatusCode >= 500 && resp.StatusCode < 600 && !responseIsStale(cacheData.CachedAt, getMaxAge(cc)+errGraceSec) {

		return genCacheResp(cacheData.Body)
	}

	return resp, err
}
