package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andream16/go-opentracing-example/src/shared/database/postgres/migrator"

	"github.com/Shopify/sarama"
	"golang.org/x/sync/errgroup"

	"github.com/andream16/go-opentracing-example/src/kafka-consumer/todo/repository"
	transportkafka "github.com/andream16/go-opentracing-example/src/kafka-consumer/transport/kafka"
	"github.com/andream16/go-opentracing-example/src/shared/database/postgres/pgxwrapper"
	"github.com/andream16/go-opentracing-example/src/shared/kafka"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
)

func main() {
	const (
		serviceName    = "kafka-consumer"
		kafkaGroupName = "kafka-consumer"
	)

	var (
		kafkaTodoTopic     string
		kafkaBrokerAddress string
		databaseDSN        string
		jaegerAgentHost    string
		jaegerAgentPort    string
	)

	for k, v := range map[string]*string{
		"KAFKA_TODO_TOPIC":     &kafkaTodoTopic,
		"KAFKA_BROKER_ADDRESS": &kafkaBrokerAddress,
		"DATABASE_DSN":         &databaseDSN,
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

	tracer := tracing.NewJaegerTracer(serviceName, jaegerAgentHost, jaegerAgentPort)
	defer tracer.Close()

	executor, err := pgxwrapper.New(ctx, databaseDSN, 10*time.Second, tracer)
	if err != nil {
		log.Fatalf("could not initialise a new executor: %v", err)
	}

	migrationCtx, migrationCancel := context.WithTimeout(ctx, 30*time.Second)
	defer migrationCancel()

	migrationConn, err := executor.GetConn(migrationCtx)
	if err != nil {
		log.Fatalf("could not get migration connection: %v", err)
	}

	defer migrationConn.Close(ctx)

	m, err := migrator.NewPgxMigrator(migrationCtx, migrationConn, "v1")
	if err != nil {
		log.Fatalf("could not create a new migration: %v", err)
	}

	m.AppendMigration(
		"create_todo_table",
		"CREATE TABLE todos (id SERIAL PRIMARY KEY, message TEXT);",
		"DROP TABLE todos;",
	)

	if err := m.Migrate(ctx); err != nil {
		log.Fatalf("could not run migration: %v", err)
	}

	repo, err := repository.New(executor)
	if err != nil {
		log.Fatalf("could not initialise a new repository: %v", err)
	}

	kafkaCfg := sarama.NewConfig()

	kafkaClient, err := kafka.NewClient([]string{kafkaBrokerAddress}, kafkaCfg, 10*time.Second)
	if err != nil {
		log.Fatalf("could not create new kafka client: %v", err)
	}

	kafkaConsumerGroup, err := kafka.NewConsumerGroup(kafkaGroupName, kafkaClient)
	if err != nil {
		log.Fatalf("could not create new kafka consumer group: %v", err)
	}

	consumer, err := transportkafka.NewConsumer(repo, tracer)
	if err != nil {
		log.Fatalf("could not create new kafka consumer: %v", err)
	}

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
