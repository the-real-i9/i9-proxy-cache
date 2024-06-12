package cacheServices

import (
	"i9pxc/appTypes"
	"net/http"
	"time"
)

/*
func getMaxAge(cc *appTypes.CacheControl, expires string) time.Duration {

	if cc.Has("s-max-age") {
		sma, _ := time.ParseDuration(cc.Get("s-max-age") + "s")

		return sma
	}

	if cc.Has("max-age") {
		sma, _ := time.ParseDuration(cc.Get("max-age") + "s")

		return sma
	}

	exp, _ := time.Parse(time.Layout, expires)

	return
} */

func responseIsStale(cachedAt time.Time, cc *appTypes.CacheControl, expires string, extraTime time.Duration, nearly bool) bool {

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

	exp, _ := time.Parse(time.Layout, expires)
	return time.Now().Before(exp.Add(-n).Add(extraTime))
}

func genCacheResp(cacheData appTypes.CacheData) (appTypes.CacheRespT, error) {
	return appTypes.CacheRespT{StatusCode: http.StatusNotModified, Header: cacheData.Header, Body: cacheData.Body}, nil
}

func filterHeader(header http.Header) http.Header {
	header = header.Clone()

	connHdVals := header.Values("Connection")
	for _, hdVal := range connHdVals {
		header.Del(hdVal)
	}

	header.Del("Connection")
	header.Del("Vary")
	header.Del("Server")

	return header
}
