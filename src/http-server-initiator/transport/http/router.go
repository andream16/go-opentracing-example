package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/andream16/go-opentracing-example/src/shared/tracing"
)

// Handler wraps a mux router.
type Handler struct {
	receiverHostname string
	httpClient       *http.Client
	router           *mux.Router
	tracer           tracing.Tracer
}

// NewHandler returns a new http handler.
func NewHandler(
	receiverHostname string,
	httpClient *http.Client,
	tracer tracing.Tracer,
) Handler {
	handler := Handler{
		httpClient:       httpClient,
		receiverHostname: receiverHostname,
		router:           mux.NewRouter(),
		tracer:           tracer,
	}

	handler.Router().HandleFunc("/initiator/todo", handler.CreateTodo).Methods(http.MethodPost)

	return handler
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
