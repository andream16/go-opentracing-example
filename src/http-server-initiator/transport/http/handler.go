package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/andream16/go-opentracing-example/src/http-server-initiator/todo"
)

func (h Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	gt, receiverURL := opentracing.GlobalTracer(), h.receiverHostname+"/receiver/todo"

	span := gt.StartSpan("initiator_todo")
	defer span.Finish()

	var t todo.Todo
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Println(fmt.Sprintf("could not deserialise request body: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(t)
	if err != nil {
		log.Println(fmt.Sprintf("could serialise todo: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, receiverURL, bytes.NewReader(b))
	if err != nil {
		log.Println(fmt.Sprintf("could not create a new http request: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, receiverURL)
	ext.HTTPMethod.Set(span, http.MethodPost)

	if err := gt.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		log.Println(fmt.Sprintf("could not inject tracing headers: %s", err))
	}

	if _, err := h.httpClient.Do(req); err != nil {
		log.Println(fmt.Sprintf("could not perform receiver request: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
