package main

import (
	"fmt"
	"i9pxc/helpers"
	"i9pxc/services/appServices"
	"i9pxc/services/cacheServices"
	"log"
	"net/http"
	"os"
)

func main() {
	err := helpers.ServerInits()
	if err != nil {
		log.Fatal(err)
	}

	cacheServerUrl := os.Getenv("CACHE_SERVER_URL")

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

		cacheKey := helpers.GenCacheKey(cacheServerUrl, r)

		cacheResp, err := cacheServices.ServeResponse(r, cacheKey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		helpers.CopyHeader(w.Header(), cacheResp.Header)

		w.WriteHeader(cacheResp.StatusCode)

		_, w_err := w.Write(cacheResp.Body)
		if w_err != nil {
			log.Println(w_err)
		}
	})

	fmt.Printf("Server listening @ %s\n", cacheServerUrl)
	http.ListenAndServe(os.Getenv("CACHE_SERVER_PORT"), nil)
}
