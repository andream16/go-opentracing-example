package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler wraps a mux router.
type Handler struct {
	receiverHostname string
	httpClient       *http.Client
	router           *mux.Router
}

// NewHandler returns a new http handler.
func NewHandler(
	receiverHostname string,
	httpClient *http.Client,
) Handler {
	handler := Handler{
		httpClient:       httpClient,
		receiverHostname: receiverHostname,
		router:           mux.NewRouter(),
	}

	handler.Router().HandleFunc("/initiator/todo", handler.CreateTodo).Methods(http.MethodPost)
	handler.Router().HandleFunc("/_health", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)

	return handler
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
