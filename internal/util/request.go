package util

import (
	"fmt"
	"io"
	"net/http"

	"github.com/antlabs/pcurl"
)

// Request request data
func Request(req *http.Request, client *http.Client) (io.Reader, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("req url %s status code %d", req.URL, resp.StatusCode)
	}
	return resp.Body, nil
}

func ParseAndRequest(url string) (*http.Request, error) {
	url = RemoveSlash(url)
	return pcurl.ParseAndRequest(url)
}
