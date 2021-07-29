package httpcli

import (
	"net/http"
)

// Built-in default HTTP client instance.
var defaultHTTPClient = NewDefaultHTTPClient()

// DefaultHTTPClient returns the built-in default HTTP client instance.
func DefaultHTTPClient() *http.Client {
	return defaultHTTPClient
}

// NewDefaultHTTPClient returns a new default HTTP client instance.
func NewDefaultHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	// Default http.Transport.MaxIdleConnsPerHost is 2.
	t.MaxIdleConnsPerHost = 10

	return &http.Client{Transport: t}
}
