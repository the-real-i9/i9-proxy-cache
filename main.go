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

		cacheRequestURL := fmt.Sprintf("%s://%s/%s", "http", r.Host, r.URL.String())

		if cacheResp, found := cacheServices.ServeRequest(r, cacheRequestURL); found {
			w.Write(cacheResp.Body)
			return
		}

		originResp, err := appServices.ForwardRequest(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		body, _ := io.ReadAll(originResp.Body)

		go func(body []byte) {
			originResp := *originResp
			originResp.Body = io.NopCloser(bytes.NewReader(body))

			cacheServices.CacheResponse(&originResp, cacheRequestURL)
		}(body)

		w.WriteHeader(originResp.StatusCode)
		w.Write(body)
	})

	fmt.Println("Server listening @ http://localhost:5000")
	http.ListenAndServe(":5000", nil)
}
