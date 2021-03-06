// Package body_size_limit the body_size_limit plugin.
package body_size_limit

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Config the plugin configuration.
// Limit a limit in bytes.
type Config struct {
	Limit int64 `json:"limit" yaml:"limit"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Limit: 52428800, // 50 MB
	}
}

// RequestSize a RequestSize plugin.
type RequestSize struct {
	next  http.Handler
	limit int64
	name  string
}

// New created a new RequestSize plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Limit <= 0 {
		return nil, fmt.Errorf("limit must be larger than zero")
	}

	return &RequestSize{
		limit: config.Limit,
		next:  next,
		name:  name,
	}, nil
}

func (a *RequestSize) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cutBody := http.MaxBytesReader(rw, req.Body, a.limit)

	all, err := ioutil.ReadAll(cutBody)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusRequestEntityTooLarge)
		return
	}

	if int64(len(all)) == a.limit {
		http.Error(rw, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(all))

	a.next.ServeHTTP(rw, req)
}
