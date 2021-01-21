// Package utils has no test file, therefore it's ignored code cover calculation.
// This is useful because utils should only contain untestable functions.
package utils

import (
	"io"
	"net/http"
)

// JSONGetWithHTTP retrieves a JSON from a URL.
func JSONGetWithHTTP(url string) (out io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err == nil {
		out = resp.Body
	}
	return out, err
}
