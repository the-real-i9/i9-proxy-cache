package appServices

import (
	"fmt"
	"net/http"
)

func ForwardRequest(r *http.Request) (*http.Response, error) {
	originRequestURL := fmt.Sprintf("%s://%s%s", "http", "localhost:8080", r.URL.String())

	req, err := http.NewRequest(r.Method, originRequestURL, r.Body)

	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)

}
