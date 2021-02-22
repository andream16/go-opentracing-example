package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	"github.com/andream16/go-opentracing-example/src/grpc-server/transport/grpc/todo"
	"github.com/andream16/go-opentracing-example/src/shared/kafka"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
)

func main() {
	const serviceName = "grpc-server"

	var (
		grpcServerPort     string
		kafkaTodoTopic     string
		kafkaBrokerAddress string
		jaegerAgentHost    string
		jaegerAgentPort    string
	)

	for k, v := range map[string]*string{
		"GRPC_SERVER_PORT":     &grpcServerPort,
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

	tracer, err := tracing.NewJaegerTracer(serviceName, jaegerAgentHost, jaegerAgentPort)
	if err != nil {
		log.Fatalf("could not create new tracer: %v", err)
	}
	defer tracer.Close()

	kafkaCfg := sarama.NewConfig()

	kafkaCfg.Producer.RequiredAcks = sarama.WaitForAll
	kafkaCfg.Producer.Retry.Max = 10
	kafkaCfg.Producer.Return.Successes = true

	kafkaClient, err := kafka.NewClient([]string{kafkaBrokerAddress}, kafkaCfg, time.Second*10)
	if err != nil {
		log.Fatalf("could not create new kafka client: %v", err)
	}

	kafkaProducer, err := kafka.NewSyncProducer(kafkaClient)
	if err != nil {
		log.Fatalf("could not create new kafka producer: %v", err)
	}

	service, err := todo.NewService(kafkaTodoTopic, kafkaProducer, tracer)
	if err != nil {
		log.Fatalf("could not create new service: %v", err)
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
	)

	todov1.RegisterTodoServiceServer(grpcSrv, service)

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
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		grpcSrv.GracefulStop()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("exiting: %v", err)
	}
}
