package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/andream16/go-opentracing-example/src/shared/health"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	internalhealth "github.com/andream16/go-opentracing-example/src/grpc-server/health"
	"github.com/andream16/go-opentracing-example/src/grpc-server/transport/grpc/todo"
	transporthttp "github.com/andream16/go-opentracing-example/src/grpc-server/transport/http"
)

func main() {
	const serviceName = "grpc-server"

	var (
		grpcServerPort     string
		kafkaTodoTopic     string
		kafkaBrokerAddress string
		jaegerAgentHost    string
		jaegerAgentPort    string
		httpServerHostname string
	)

	for k, v := range map[string]*string{
		"GRPC_SERVER_PORT":     &grpcServerPort,
		"HTTP_SERVER_HOSTNAME": &httpServerHostname,
		"KAFKA_TODO_TOPIC":     &kafkaTodoTopic,
		"KAFKA_BROKER_ADDRESS": &kafkaBrokerAddress,
		"JAEGER_AGENT_HOST":    &jaegerAgentHost,
		"JAEGER_AGENT_PORT":    &jaegerAgentPort,
	} {
		var ok bool
		*v, ok = os.LookupEnv(k)
		if !ok {
			log.Fatalf("missing environment variable %s", k)
		}
	}

	var ctx, cancel = context.WithCancel(context.Background())
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

	kafkaCfg := sarama.NewConfig()

	kafkaCfg.Producer.RequiredAcks = sarama.WaitForAll
	kafkaCfg.Producer.Retry.Max = 10
	kafkaCfg.Producer.Return.Successes = true

	kafkaProducer, err := sarama.NewSyncProducer([]string{kafkaBrokerAddress}, kafkaCfg)
	if err != nil {
		log.Fatalf("could not create new kafka producer: %v", err)
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
	)

	todov1.RegisterTodoServiceServer(grpcSrv, todo.NewService(kafkaTodoTopic, kafkaProducer))

	handler := transporthttp.NewHandler(
		health.NewManager(internalhealth.NewKafkaChecker(nil)),
	)

	httpServ := &http.Server{
		Addr:         httpServerHostname,
		Handler:      handler.Router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		l, err := net.Listen("tcp", ":"+grpcServerPort)
		if err != nil {
			return fmt.Errorf("could prepare for grpc dialing: %w", err)
		}

		defer l.Close()

		log.Println(fmt.Sprintf("serving traffic at 0.0.0.0:%s ...", grpcServerPort))

		return grpcSrv.Serve(l)
	})

	g.Go(func() error {
		log.Println(fmt.Sprintf("serving traffic at 0.0.0.0:%s ...", httpServerHostname))
		return httpServ.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		_ = httpServ.Shutdown(ctx)
		grpcSrv.GracefulStop()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("exiting: %v", err)
	}
}
