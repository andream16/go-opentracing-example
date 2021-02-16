package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler wraps a mux router.
type Handler struct {
	httpClient       *http.Client
	receiverHostname string
	router           *mux.Router
}

// NewHandler returns a new http handler.
func NewHandler(httpClient *http.Client, receiverHostname string) Handler {
	handler := Handler{
		httpClient:       httpClient,
		receiverHostname: receiverHostname,
		router:           mux.NewRouter(),
	}

	handler.Router().HandleFunc("/initiator/todo", handler.CreateTodo).Methods(http.MethodPost)
	handler.Router().HandleFunc("/_health", func(w http.ResponseWriter, r *http.Request) {})

	return handler
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
