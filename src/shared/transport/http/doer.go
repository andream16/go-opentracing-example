package http

import "net/http"

// Doer is an interface used for testing.
type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}
