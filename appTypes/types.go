package appTypes

import (
	"net/http"
	"strings"
	"time"
)

type CacheRespT struct {
	StatusCode int
	Body       []byte
}

type CacheData struct {
	Header   http.Header
	Trailer  http.Header
	Body     []byte
	CachedAt time.Time
}

type CacheControl struct {
	dirsMap map[string]any
}

func (cc *CacheControl) Parse(ccDirs []string) {

	dirsMap := make(map[string]any)

	for _, ccd := range ccDirs {
		key, value, _ := strings.Cut(ccd, "=")

		dirsMap[key] = strings.ToLower(value)
	}

	cc.dirsMap = dirsMap
}

func (cc CacheControl) Has(key string) bool {
	return cc.dirsMap[key] != nil
}

func (cc CacheControl) Get(key string) string {
	if cc.dirsMap[key] == nil {
		panic("check if key exists before getting")
	}

	return cc.dirsMap[key].(string)
}
