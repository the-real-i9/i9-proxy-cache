package cacheServices

import (
	"fmt"
	"i9pxc/appTypes"
	"net/http"
	"time"
)

func getMaxAge(cc appTypes.CacheControl) float64 {
	var maxAge float64

	if cc.Has("s-max-age") {
		fmt.Sscanf(cc.Get("s-max-age"), "%f", &maxAge)
		return maxAge
	}

	fmt.Sscanf(cc.Get("max-age"), "%f", &maxAge)

	return maxAge

}

func responseIsStale(cachedAt time.Time, maxAge float64) bool {

	return time.Since(cachedAt).Seconds() > maxAge
}

func responseIsNearlyStale(cachedAt time.Time, maxAge float64) bool {

	return time.Since(cachedAt).Seconds() > (0.9 * maxAge)
}

func genCacheResp(body []byte) (appTypes.CacheRespT, error) {
	return appTypes.CacheRespT{StatusCode: http.StatusOK, Body: body}, nil
}
