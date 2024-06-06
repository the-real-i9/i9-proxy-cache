package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// url key
		originHost := "http://localhost:8080"
		r.URL.Host = originHost

		reqUrl := r.URL.String()

		// the response to cache
		proxyRes, err := http.Get(reqUrl)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}

		proxyRes.Write(w)
	})

	http.ListenAndServe(":5000", nil)
	fmt.Println("Server listening @ http://localhost:5000")
}
