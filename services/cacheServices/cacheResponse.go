package cacheServices

import (
	"i9pxc/appTypes"
	"net/http"
)

func respIsCacheable(resp *http.Response) bool {
	cc := appTypes.CacheCtrl{}
	cc.Parse(resp.Header.Values("Cache-Control"))

	if resp.StatusCode != http.StatusOK {
		return false
	}

	if cc.Has("no-store") && !(cc.Has("must-understand") && resp.StatusCode == http.StatusOK) {

		return false
	}

	if cc.Has("private") {
		return false
	}

	if resp.Header.Get("Authorization") != "" && !(cc.Has("public") || cc.Has("s-max-age") || cc.Has("must-revalidate")) {

		return false
	}

	if !(cc.Has("public") || cc.Has("max-age") || cc.Has("s-max-age") || cc.Has("no-cache") || resp.Header.Get("Expires") != "") {
		return false
	}

	return true
}

func filterHeader(header http.Header) http.Header {
	connHdVals := header.Values("Connection")
	for _, hdVal := range connHdVals {
		header.Del(hdVal)
	}

	header.Del("Connection")

	return header
}

func CacheResponse(originResp *http.Response, cacheRequestURL string) {

	if !respIsCacheable(originResp) {
		return
	}

	header := filterHeader(originResp.Header)
}
