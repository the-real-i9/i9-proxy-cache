package main

import (
	"context"
	"encoding/json"
	"fmt"
	"i9pxc/db"
	"i9pxc/helpers"
	"io"
	"net/http"
	"testing"
)

func TesResp(t *testing.T) {
	err := helpers.ServerInits()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get("http://google.com")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	bodyData, _ := io.ReadAll(resp.Body)

	cacheData := map[string]any{
		"header": resp.Header,
		"body":   bodyData,
	}

	cacheDataBt, err := json.Marshal(cacheData)
	if err != nil {
		t.Fatal(err)
	}

	rdb_err := db.RedisDB.Set(context.Background(), "http://google.com", cacheDataBt, 0).Err()
	if rdb_err != nil {
		t.Fatal(rdb_err)
	}

	var hd struct {
		Header http.Header
		Body   []byte
	}

	res, _ := db.RedisDB.Get(context.Background(), "http://google.com").Result()

	u_err := json.Unmarshal([]byte(res), &hd)
	if u_err != nil {
		t.Error(u_err)
	}

	fmt.Println(hd.Header.Get("Connection"))
}
