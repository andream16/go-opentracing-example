package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-receiver/transport/http"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tracer, err := tracing.NewJaegerTracer(serviceName, jaegerAgentHost, jaegerAgentPort)
	if err != nil {
		log.Fatalf("could not create new tracer: %v", err)
	}
	defer tracer.Close()

	conn, err := grpc.Dial(
		grpcServerHostname,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		log.Fatalf("could not initialise grpc client: %v", err)
	}

	handler, err := transporthttp.NewHandler(todov1.NewTodoServiceClient(conn), tracer)
	if err != nil {
		log.Fatalf("could not create a new handler: %v", err)
	}

	server := &http.Server{
		Addr:         httpServerHostname,
		Handler:      handler.Router(),
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
