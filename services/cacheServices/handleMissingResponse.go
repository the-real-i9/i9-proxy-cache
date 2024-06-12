package cacheServices

import (
	"bytes"
	"i9pxc/appTypes"
	"i9pxc/services/appServices"
	"io"
	"net/http"
)

func handleMissingResp(r *http.Request, cacheRequestKey string) (appTypes.CacheRespT, error) {
	originResp, err := appServices.ForwardRequest(r)
	if err != nil {
		return appTypes.CacheRespT{}, err
	}

	defer originResp.Body.Close()

	body, _ := io.ReadAll(originResp.Body)

	go func(body []byte) {
		originResp := *originResp
		originResp.Body = io.NopCloser(bytes.NewReader(body))

		CacheResponse(&originResp, cacheRequestKey)
	}(body)

	return appTypes.CacheRespT{StatusCode: originResp.StatusCode, Body: body}, nil
}
