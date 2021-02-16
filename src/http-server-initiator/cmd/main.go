package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"golang.org/x/sync/errgroup"

	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-initiator/transport/http"
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

	handler := transporthttp.NewHandler(httpClient, httpServerReceiverHostname)

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
