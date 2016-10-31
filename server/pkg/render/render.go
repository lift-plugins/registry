package render

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// Option allows setting rendering options.
type Option func(*options)
type options struct {
	status int
	cache  bool
	body   interface{}
}

// WithStatus sets the response HTTP status code.
func WithStatus(status int) Option {
	return func(o *options) {
		o.status = status
	}
}

// WithBody sets the response body.
func WithBody(body interface{}) Option {
	return func(o *options) {
		o.body = body
	}
}

// WithCache enables caching.
func WithCache() Option {
	return func(o *options) {
		o.cache = true
	}
}

// JSON renders JSON content and sends it to the HTTP client. It supports caching.
func JSON(w http.ResponseWriter, opts ...Option) error {
	if &w == nil {
		return fmt.Errorf("You must provide a valid http.ResponseWriter")
	}

	options := &options{
		cache:  false,
		status: http.StatusOK,
	}

	for _, opt := range opts {
		opt(options)
	}

	headers := w.Header()
	headers.Set("Content-Type", "application/json; charset=utf-8")

	if !options.cache {
		headers.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		headers.Set("Pragma", "no-cache")
		headers.Set("Expires", "0")
	}

	jsonBytes, err := json.Marshal(options.body)
	if err != nil {
		return errors.Wrap(err, "failed marshaling JSON response body")
	}

	headers.Set("Content-Length", strconv.Itoa(len(jsonBytes)))
	w.WriteHeader(options.status)

	if _, err := w.Write(jsonBytes); err != nil {
		return errors.Wrap(err, "failed writing JSON response body")
	}

	return nil
}
