package main

import (
	"bytes"
	"fmt"
	"i9pxc/helpers"
	"i9pxc/services/appServices"
	"i9pxc/services/cacheServices"
	"io"
	"log"
	"net/http"
)

func main() {
	err := helpers.ServerInits()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			resp, err := appServices.ForwardRequest(r)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

			resp.Write(w)
			return
		}

		if cacheResp, found := cacheServices.ServeRequest(r); found {
			w.Write(cacheResp.Body)
			return
		}

		resp, err := appServices.ForwardRequest(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		body, _ := io.ReadAll(resp.Body)

		go func(body []byte) {
			resp := *resp
			resp.Body = io.NopCloser(bytes.NewReader(body))

			cacheServices.CacheResponse(&resp)
		}(body)

		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	fmt.Println("Server listening @ http://localhost:5000")
	http.ListenAndServe(":5000", nil)
}
