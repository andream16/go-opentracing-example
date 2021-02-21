package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/andream16/go-opentracing-example/src/shared/tracing"
	transporthttp "github.com/andream16/go-opentracing-example/src/shared/transport/http"
)

// Handler wraps a mux router.
type Handler struct {
	receiverHostname string
	doer             transporthttp.Doer
	router           *mux.Router
	tracer           tracing.Tracer
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
func NewHandler(
	receiverHostname string,
	doer transporthttp.Doer,
	tracer tracing.Tracer,
) (Handler, error) {
	handler := Handler{}

	switch {
	case receiverHostname == "":
		return handler, InvalidHandlerParameterError{parameter: "receiverHostname", reason: "cannot be empty"}
	case doer == nil:
		return handler, InvalidHandlerParameterError{parameter: "doer", reason: "cannot be nil"}
	case tracer == nil:
		return handler, InvalidHandlerParameterError{parameter: "tracer", reason: "cannot be nil"}
	}

	handler.doer = doer
	handler.receiverHostname = receiverHostname
	handler.router = mux.NewRouter()
	handler.tracer = tracer

	handler.Router().HandleFunc("/initiator/todo", handler.CreateTodo).Methods(http.MethodPost)

	return handler, nil
}

// Router returns the inner router.
func (h Handler) Router() *mux.Router {
	return h.router
}
