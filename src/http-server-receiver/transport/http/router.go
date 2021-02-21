package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
)

// Handler wraps a mux router.
type Handler struct {
	todoSvcClient todov1.TodoServiceClient
	router        *mux.Router
	tracer        opentracing.Tracer
}

// InvalidHandlerParameterError is used when an invalid parameter is passed to NewHandler.
type InvalidHandlerParameterError struct {
	parameter string
	reason    string
}

func (i InvalidHandlerParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s: %s", i.parameter, i.reason)
}

// NewHandler returns a new http handler.
func NewHandler(todoSvcClient todov1.TodoServiceClient, tracer opentracing.Tracer) (Handler, error) {
	handler := Handler{}

	switch {
	case todoSvcClient == nil:
		return handler, InvalidHandlerParameterError{parameter: "todoSvcClient", reason: "cannot be nil"}
	case tracer == nil:
		return handler, InvalidHandlerParameterError{parameter: "tracer", reason: "cannot be nil"}
	}

	handler.todoSvcClient = todoSvcClient
	handler.tracer = tracer
	handler.router = mux.NewRouter()

	handler.Router().HandleFunc("/receiver/todo", handler.CreateTodo).Methods(http.MethodPost)

	return handler, nil
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
