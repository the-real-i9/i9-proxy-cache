package cacheServices

import (
	"i9pxc/appTypes"
	"net/http"
)

func ServeResponse(r *http.Request, cacheRequestKey string) (appTypes.CacheRespT, error) {
	cacheData, found := retrieveResp(cacheRequestKey)
	if !found {
		return handleMissingResp(r, cacheRequestKey)
	}

	cc := appTypes.CacheControl{}
	cc.Parse(cacheData.Header.Values("Cache-Control"))

	if cc.Has("no-cache") {
		return strictRevalidate(r, cacheData, cacheRequestKey)
	}

	if responseIsStale(cacheData.CachedAt, getMaxAge(cc)) {
		if cc.Has("must-revalidate") {
			return strictRevalidate(r, cacheData, cacheRequestKey)
		}

		if cc.Has("stale-while-revalidate") {
			return staleWhileRevalidate(r, cc, cacheData, cacheRequestKey)
		}

		if cc.Has("stale-if-error") {
			return staleIfError(r, cc, cacheData, cacheRequestKey)
		}

		return nonStrictRevalidate(r, cacheData, cacheRequestKey)
	}

	if responseIsNearlyStale(cacheData.CachedAt, getMaxAge(cc)) {
		go revalidate(r, cacheData, cacheRequestKey)
	}

	return genCacheResp(cacheData)
}
