package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-initiator/transport/http"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
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
		httpClient  = &http.Client{Timeout: 10 * time.Second}
	)

	defer cancel()

	tracer, err := tracing.NewJaegerTracer(serviceName, jaegerAgentHost, jaegerAgentPort)
	if err != nil {
		log.Fatalf("could not create new tracer: %v", err)
	}
	defer tracer.Close()

	handler, err := transporthttp.NewHandler(
		httpServerReceiverHostname,
		httpClient,
		tracer,
	)
	if err != nil {
		log.Fatalf("could not create handler: %v", err)
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
