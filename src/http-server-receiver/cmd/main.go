package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	"github.com/andream16/go-opentracing-example/src/http-server-receiver/todo"
)

func main() {
	const serviceName = "http-server-receiver"

	var (
		httpServerHostname string
		grpcServerHostname string
		jaegerAgentHost    string
		jaegerAgentPort    string
	)

	for k, v := range map[string]*string{
		"HTTP_SERVER_HOSTNAME": &httpServerHostname,
		"GRPC_SERVER_HOSTNAME": &grpcServerHostname,
		"JAEGER_AGENT_HOST":    &jaegerAgentHost,
		"JAEGER_AGENT_PORT":    &jaegerAgentPort,
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

	conn, err := grpc.Dial(
		grpcServerHostname,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(grpc_opentracing.StreamClientInterceptor())),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(grpc_opentracing.UnaryClientInterceptor())),
	)
	if err != nil {
		log.Fatalf("could not initialise grpc client: %v", err)
	}

	todoSvcClient := todov1.NewTodoServiceClient(conn)

	router.HandleFunc("/receiver/todo", func(w http.ResponseWriter, r *http.Request) {
		gt := opentracing.GlobalTracer()

		spanCtx, err := gt.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			log.Println(fmt.Sprintf("could not extract tracing headers: %s", err))
		}

		span := gt.StartSpan("receiver_todo", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		var t todo.Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			log.Println(fmt.Sprintf("could not deserialise request body: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		if _, err := todoSvcClient.Create(r.Context(), &todov1.CreateRequest{Message: t.Message}); err != nil {
			log.Println(fmt.Sprintf("could not create todo: %s", err))
			w.WriteHeader(http.StatusServiceUnavailable)
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
