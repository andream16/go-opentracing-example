package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"golang.org/x/sync/errgroup"

	"github.com/andream16/go-opentracing-example/src/kafka-consumer/transport/kafka"
)

func main() {
	const (
		serviceName    = "kafka-consumer"
		kafkaGroupName = "kafka-consumer"
	)

	var (
		kafkaTodoTopic     string
		kafkaBrokerAddress string
		jaegerAgentHost    string
		jaegerAgentPort    string
	)

	for k, v := range map[string]*string{
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

	kafkaConsumerGroup, err := sarama.NewConsumerGroup(
		[]string{kafkaBrokerAddress},
		kafkaGroupName,
		kafkaCfg,
	)
	if err != nil {
		log.Fatalf("could not create new kafka consumer group: %v", err)
	}

	consumer := kafka.NewConsumer()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			if err := kafkaConsumerGroup.Consume(
				ctx,
				[]string{kafkaTodoTopic},
				consumer,
			); err != nil {
				cancel()
				return fmt.Errorf("received fatal error from consumer: %w", err)
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case err := <-kafkaConsumerGroup.Errors():
				log.Println(fmt.Sprintf("received error while consuming: %v", err))
			default:
			}
		}
	})

	g.Go(func() error {
		<-ctx.Done()
		return kafkaConsumerGroup.Close()
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("exiting: %v", err)
	}
}
