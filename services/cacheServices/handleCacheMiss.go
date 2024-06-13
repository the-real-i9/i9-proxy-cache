package cacheServices

import (
	"i9pxc/appTypes"
	"i9pxc/helpers"
	"i9pxc/services/appServices"
	"io"
	"net/http"
)

func handleCacheMiss(r *http.Request, cacheKey string) (appTypes.CacheRespT, error) {
	originResp, err := appServices.ForwardRequest(r)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	defer originResp.Body.Close()

	body, _ := io.ReadAll(originResp.Body)

	go CacheResponse(originResp, cacheKey, body)

	return appTypes.CacheRespT{StatusCode: originResp.StatusCode, Header: helpers.FilterHeader(originResp.Header), Body: body}, nil
}
