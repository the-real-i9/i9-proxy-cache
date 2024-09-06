package appServices

import (
	"fmt"
	"net/http"
	"os"
)

func ForwardRequest(r *http.Request) (*http.Response, error) {
	defer r.Body.Close()

	originRequestURL := fmt.Sprintf("%s%s", os.Getenv("ORIGIN_SERVER_URL"), r.URL.String())

	req, err := http.NewRequest(r.Method, originRequestURL, r.Body)

	req.Header = r.Header.Clone()

	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func RevalidationRequest(reqURL, lastModified, eTag string) (*http.Response, error) {
	originRequestURL := fmt.Sprintf("%s%s", os.Getenv("ORIGIN_SERVER_URL"), reqURL)

	req, err := http.NewRequest("GET", originRequestURL, nil)
	if err != nil {
		return nil, err
	}

	if eTag != "" {
		req.Header.Set("If-None-Match", eTag)
	} else if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}

	return http.DefaultClient.Do(req)
}
