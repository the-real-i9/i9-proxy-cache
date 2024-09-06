package cacheServices

import (
	"i9pxc/src/appTypes"
	"i9pxc/src/helpers"
	"i9pxc/src/services/appServices"
	"io"
	"net/http"
	"time"
)

func revalidate(r *http.Request, cacheData appTypes.CacheData, cacheKey string) (appTypes.CacheRespT, error) {
	resp, err := appServices.RevalidationRequest(r.URL.String(), cacheData.Header.Get("Last-Modified"), cacheData.Header.Get("ETag"))
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	if resp.StatusCode == http.StatusNotModified {
		go RefreshCacheResponse(cacheKey)
		return helpers.GenCacheResp(cacheData)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	go CacheResponse(resp, cacheKey, body)

	return appTypes.CacheRespT{StatusCode: resp.StatusCode, Header: helpers.FilterHeader(resp.Header), Body: body}, nil
}

func nonStrictRevalidate(r *http.Request, cacheData appTypes.CacheData, cacheKey string) (appTypes.CacheRespT, error) {
	resp, err := revalidate(r, cacheData, cacheKey)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return helpers.GenCacheResp(cacheData)
	}

	return resp, err
}

func strictRevalidate(r *http.Request, cacheData appTypes.CacheData, cacheKey string) (appTypes.CacheRespT, error) {
	return revalidate(r, cacheData, cacheKey)
}

func staleWhileRevalidate(r *http.Request, cc *appTypes.CacheControl, cacheData appTypes.CacheData, cacheKey string) (appTypes.CacheRespT, error) {
	extraTime, _ := time.ParseDuration(cc.Get("stale-while-revalidate") + "s")

	if helpers.ResponseIsStaleWithExtraTime(cacheData.CachedAt, cc, cacheData.Header.Get("Expires"), extraTime) {
		return revalidate(r, cacheData, cacheKey)
	}

	return helpers.GenCacheResp(cacheData)
}

func staleIfError(r *http.Request, cc *appTypes.CacheControl, cacheData appTypes.CacheData, cacheKey string) (appTypes.CacheRespT, error) {
	resp, err := revalidate(r, cacheData, cacheKey)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	extraTime, _ := time.ParseDuration(cc.Get("stale-if-error") + "s")

	if resp.StatusCode >= 500 && resp.StatusCode < 600 && !helpers.ResponseIsStaleWithExtraTime(cacheData.CachedAt, cc, cacheData.Header.Get("Expires"), extraTime) {
		return helpers.GenCacheResp(cacheData)
	}

	return resp, err
}
