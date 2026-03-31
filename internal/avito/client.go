// Package avito provides Avito API clients.
package avito

import (
	"net/http"
	"time"
)

// Client wraps shared Avito HTTP dependencies.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient returns a new API client.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// NewHTTPClient returns a default HTTP client for Avito API.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
