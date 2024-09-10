package cacheServices

import (
	"i9pxc/appTypes"
	"i9pxc/helpers"
	"net/http"
)

func handleCacheHit(r *http.Request, cacheData appTypes.CacheData, cacheKey string) (appTypes.CacheRespT, error) {
	cc := &appTypes.CacheControl{}
	cc.Parse(cacheData.Header.Get("Cache-Control"))

	if cc.Has("no-cache") {
		return strictRevalidate(r, cacheData, cacheKey)
	}

	if helpers.ResponseIsStale(cacheData.CachedAt, cc, cacheData.Header.Get("Expires")) {
		if cc.Has("must-revalidate") {
			return strictRevalidate(r, cacheData, cacheKey)
		}

		if cc.Has("stale-while-revalidate") {
			return staleWhileRevalidate(r, cc, cacheData, cacheKey)
		}

		if cc.Has("stale-if-error") {
			return staleIfError(r, cc, cacheData, cacheKey)
		}

		return nonStrictRevalidate(r, cacheData, cacheKey)
	}

	if helpers.ResponseIsNearlyStale(cacheData.CachedAt, cc, cacheData.Header.Get("Expires")) {
		go revalidate(r, cacheData, cacheKey)
	}

	return helpers.GenCacheResp(cacheData)
}
