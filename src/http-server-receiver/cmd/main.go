package main

import (
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

	"github.com/andream16/go-opentracing-example/src/http-server-receiver/todo"
)

func main() {
	const (
		hostname    = "0.0.0.0:8081"
		serviceName = "http-server-receiver"
	)

	var (
		ctx, cancel = context.WithCancel(context.Background())
		router      = mux.NewRouter()
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

	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	router.HandleFunc("/receiver/todo", func(w http.ResponseWriter, r *http.Request) {
		trc := opentracing.GlobalTracer()

		spanCtx, err := trc.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			log.Println(fmt.Sprintf("could not extract tracing headers: %s", err))
		}

		span := trc.StartSpan("receiver_todo", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		var t todo.Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			log.Println(fmt.Sprintf("could not deserialise request body: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		log.Println(fmt.Sprintf("new todo with messageL: %s", t.Message))
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
