package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/andream16/go-opentracing-example/src/shared/health"
)

// Handler wraps a mux router.
type Handler struct {
	router *mux.Router
}

// NewHandler returns a new http handler.
func NewHandler(checker health.Manager) Handler {
	handler := Handler{
		router: mux.NewRouter(),
	}

	handler.Router().HandleFunc("/_health", health.Handler(checker)).Methods(http.MethodGet)

	return handler
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
