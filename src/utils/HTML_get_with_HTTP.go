// Package utils has no test file, therefore it's ignored code cover calculation.
// This is useful because utils should only contain untestable functions.
package utils

import (
	"io"
	"net/http"
)

// HTMLGetWithHTTP retrieves a HTML page from a URL.
func HTMLGetWithHTTP(url string) (redirect string, out io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err == nil {
		out = resp.Body
		redirect = resp.Request.URL.String()
	}
	return redirect, out, err
}
