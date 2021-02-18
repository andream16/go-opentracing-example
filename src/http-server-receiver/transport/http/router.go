package http

import (
	"net/http"

	"github.com/gorilla/mux"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
)

// Handler wraps a mux router.
type Handler struct {
	todoSvcClient todov1.TodoServiceClient
	router        *mux.Router
}

// NewHandler returns a new http handler.
func NewHandler(todoSvcClient todov1.TodoServiceClient) Handler {
	handler := Handler{
		todoSvcClient: todoSvcClient,
		router:        mux.NewRouter(),
	}

	handler.Router().HandleFunc("/receiver/todo", handler.CreateTodo).Methods(http.MethodPost)
	handler.Router().HandleFunc("/_health", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)

	return handler
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
