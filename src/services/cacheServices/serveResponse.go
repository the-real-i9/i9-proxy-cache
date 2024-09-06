package cacheServices

import (
	"i9pxc/src/appTypes"

	"net/http"
)

func ServeResponse(r *http.Request, cacheKey string) (appTypes.CacheRespT, error) {
	cacheData, found := retrieveResp(cacheKey)
	if !found {
		return handleCacheMiss(r, cacheKey)
	}

	return handleCacheHit(r, cacheData, cacheKey)
}
