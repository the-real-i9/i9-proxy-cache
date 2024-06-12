package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestTry(t *testing.T) {
	/* err := helpers.ServerInits()
	if err != nil {
		t.Fatal(err)
	} */

	res, _ := http.Get("https://github.com")

	res1 := *res

	res1.Request = nil

	bt, _ := json.MarshalIndent(res1, "", " ")

	os.Stdout.Write(bt)
}
