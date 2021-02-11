package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"golang.org/x/sync/errgroup"

	"github.com/andream16/go-opentracing-example/src/http-server-initiator/todo"
)

func main() {
	const (
		hostname         = "0.0.0.0:8080"
		serviceName      = "http-server-initiator"
		receiverHostname = "http://http-server-receiver:8081"
	)

	var (
		ctx, cancel = context.WithCancel(context.Background())
		router      = mux.NewRouter()
		httpClient  = &http.Client{Timeout: 10 * time.Second}
	)

	defer cancel()

	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaegerlog.StdLogger))
	if err != nil {
		log.Fatalf("could not initialise tracer: %v", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	router.HandleFunc("/initiator/todo", func(w http.ResponseWriter, r *http.Request) {
		const receiverURL = receiverHostname + "/receiver/todo"

		span := opentracing.GlobalTracer().StartSpan("initiator_todo")
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, receiverHostname+"/receiver/todo", bytes.NewReader(b))
		if err != nil {
			log.Println(fmt.Sprintf("could not create a new http request: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, receiverURL)
		ext.HTTPMethod.Set(span, http.MethodPost)

		if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
			log.Println(fmt.Sprintf("could not inject tracing headers: %s", err))
		}

		if _, err := httpClient.Do(req); err != nil {
			log.Println(fmt.Sprintf("could not perform receiver request: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodPost)

	server := &http.Server{
		Addr:         hostname,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Println(fmt.Sprintf("serving traffic at %s ...", hostname))
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("error while serving: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("exiting: %v", err)
	}
}
