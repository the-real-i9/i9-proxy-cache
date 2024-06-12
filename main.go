package main

import (
	"fmt"
	"i9pxc/helpers"
	"i9pxc/services/appServices"
	"i9pxc/services/cacheServices"
	"log"
	"net/http"
	"os"
	"slices"
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

		vary := r.Header.Values("Vary")
		slices.Sort(vary)

		cacheRequestKey := fmt.Sprintf("%s%s ~ %s", cacheServerUrl, r.URL.String(), vary)

		cacheResp, err := cacheServices.ServeResponse(r, cacheRequestKey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		for key, values := range cacheResp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(cacheResp.StatusCode)

		w.Write(cacheResp.Body)
	})

	fmt.Printf("Server listening @ %s\n", cacheServerUrl)
	http.ListenAndServe(os.Getenv("CACHE_SERVER_PORT"), nil)
}
