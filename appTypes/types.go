package appTypes

import (
	"net/http"
	"strings"
)

type CacheResp struct {
	Header  http.Header
	Trailer http.Header
	Body    []byte
}

type CacheCtrl struct {
	dirsMap map[string]any
}

func (dm *CacheCtrl) Parse(cacheCtrlDirs []string) {

	dirsMap := make(map[string]any)

	for _, ccd := range cacheCtrlDirs {
		key, value, _ := strings.Cut(ccd, "=")

		dirsMap[key] = strings.ToLower(value)
	}

	dm.dirsMap = dirsMap
}

func (dm CacheCtrl) Has(key string) bool {
	return dm.dirsMap[key] != nil
}

func (dm CacheCtrl) Get(key string) string {
	if dm.dirsMap[key] == nil {
		panic("check if key exists before getting")
	}

	return dm.dirsMap[key].(string)
}
