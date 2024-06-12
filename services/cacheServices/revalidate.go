package cacheServices

import (
	"i9pxc/appTypes"
	"i9pxc/services/appServices"
	"io"
	"net/http"
)

func revalidate(r *http.Request, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	resp, err := appServices.RevalidateRequest(r.URL.String(), cacheData.Header.Get("Last-Modified"), cacheData.Header.Get("ETag"))
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	if resp.StatusCode == http.StatusNotModified {
		go RefreshCacheResponse(cacheRequestKey)
		return genCacheResp(cacheData)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	go CacheResponse(resp, cacheRequestKey, body)

	return appTypes.CacheRespT{StatusCode: resp.StatusCode, Header: filterHeader(resp.Header), Body: body}, nil
}

func nonStrictRevalidate(r *http.Request, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	resp, err := revalidate(r, cacheData, cacheRequestKey)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return genCacheResp(cacheData)
	}

	return resp, err
}

func strictRevalidate(r *http.Request, cacheData appTypes.CacheData, cacheRequestKey string) (appTypes.CacheRespT, error) {
	return revalidate(r, cacheData, cacheRequestKey)
}
