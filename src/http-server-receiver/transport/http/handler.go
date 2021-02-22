package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	"github.com/andream16/go-opentracing-example/src/shared/todo"
)

func (h Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	spanCtx, err := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Println(fmt.Sprintf("could not extract tracing headers: %s", err))
	}

	span := h.tracer.StartSpan("receiver_todo", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	var t todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Println(fmt.Sprintf("could not deserialise request body: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if _, err := h.todoSvcClient.Create(
		opentracing.ContextWithSpan(r.Context(), span),
		&todov1.CreateRequest{Message: t.Message},
	); err != nil {
		log.Println(fmt.Sprintf("could not create todo: %s", err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}
