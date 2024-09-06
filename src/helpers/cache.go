package helpers

import (
	"i9pxc/src/appTypes"
	"log"
	"net/http"
	"strings"
	"time"
)

func isExp(cachedAt time.Time, cc *appTypes.CacheControl, expiresHeader string, extraTime time.Duration, nearly bool) bool {

	n := 0 * time.Second

	if nearly {
		n = 24 * time.Hour
	}

	if cc.Has("s-max-age") {
		sma, _ := time.ParseDuration(cc.Get("s-max-age") + "s")

		return time.Since(cachedAt) > ((sma - n) + extraTime)
	}

	if cc.Has("max-age") {
		ma, _ := time.ParseDuration(cc.Get("max-age") + "s")

		return time.Since(cachedAt) > ((ma - n) + extraTime)
	}

	if expiresHeader == "" {
		log.Panicln("helpers: cache.go: isExp: none of s-max-age, max-age, and Expires header is provided. This response should't be in the cache, debug your implementation")
	}

	exp, err := time.Parse(time.Layout, expiresHeader)
	if err != nil {
		log.Panicln("helpers: cache.go: isExp:", err)
	}

	return time.Now().Before(exp.Add(-n).Add(extraTime))
}

func ResponseIsStale(cachedAt time.Time, cc *appTypes.CacheControl, expiresHeader string) bool {
	return isExp(cachedAt, cc, expiresHeader, 0, false)
}

func ResponseIsNearlyStale(cachedAt time.Time, cc *appTypes.CacheControl, expiresHeader string) bool {
	return isExp(cachedAt, cc, expiresHeader, 0, true)
}

func ResponseIsStaleWithExtraTime(cachedAt time.Time, cc *appTypes.CacheControl, expiresHeader string, extraTime time.Duration) bool {
	return isExp(cachedAt, cc, expiresHeader, extraTime, false)
}

func GenCacheResp(cacheData appTypes.CacheData) (appTypes.CacheRespT, error) {
	return appTypes.CacheRespT{StatusCode: http.StatusOK, Header: cacheData.Header, Body: cacheData.Body}, nil
}

func FilterHeader(header http.Header) http.Header {
	header = header.Clone()

	hopByHopHeaders := strings.Split(header.Get("Connection"), ", ")
	for _, hd := range hopByHopHeaders {
		header.Del(hd)
	}

	header.Del("Connection")
	header.Del("Vary")
	header.Set("Server", "Go-http-server")

	return header
}

func GenCacheKey(cacheServerUrl string, r *http.Request) string {
	vary := r.Header.Get("Vary")
	if vary != "" {
		vary = " ~ " + vary
	}

	return cacheServerUrl + r.URL.String() + vary
}
