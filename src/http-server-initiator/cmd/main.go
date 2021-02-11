package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"golang.org/x/sync/errgroup"

	"github.com/andream16/go-opentracing-example/src/http-server-initiator/todo"
)

func main() {
	const serviceName = "http-server-initiator"

	var (
		httpServerHostname         string
		jaegerAgentHost            string
		jaegerAgentPort            string
		httpServerReceiverHostname string
	)

	for k, v := range map[string]*string{
		"HTTP_SERVER_HOSTNAME":          &httpServerHostname,
		"JAEGER_AGENT_HOST":             &jaegerAgentHost,
		"JAEGER_AGENT_PORT":             &jaegerAgentPort,
		"HTTP_SERVER_RECEIVER_HOSTNAME": &httpServerReceiverHostname,
	} {
		var ok bool
		*v, ok = os.LookupEnv(k)
		if !ok {
			log.Fatalf("missing environment variable %s", k)
		}
	}

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
			LogSpans:           true,
			LocalAgentHostPort: jaegerAgentHost + ":" + jaegerAgentPort,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerlog.StdLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatalf("could not initialise tracer: %v", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	router.HandleFunc("/initiator/todo", func(w http.ResponseWriter, r *http.Request) {
		receiverURL := httpServerReceiverHostname + "/receiver/todo"

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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, httpServerReceiverHostname+"/receiver/todo", bytes.NewReader(b))
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
		Addr:         httpServerHostname,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Println(fmt.Sprintf("serving traffic at %s ...", httpServerHostname))
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("error while serving: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("exiting: %v", err)
	}
}
