package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	"github.com/andream16/go-opentracing-example/src/kafka-consumer/todo/repository"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
)

const spanName = "todo_consumer"

// Consumer represent a kafka transport consumer.
type Consumer struct {
	repo   repository.Creator
	tracer tracing.Tracer
}

// NewConsumer returns a new consumer.
func NewConsumer(repo repository.Creator, tracer tracing.Tracer) (Consumer, error) {
	switch {
	case repo == nil:
		return Consumer{}, errors.New("repo could not be nil")
	case tracer == nil:
		return Consumer{}, errors.New("tracer could not be nil")
	}
	return Consumer{
		repo:   repo,
		tracer: tracer,
	}, nil
}

func (c Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := c.ReceivedMessage(message); err != nil {
			log.Printf("could not create todo, skipping message: %v", err)
		}
		session.MarkMessage(message, "")
	}

	return nil
}

// ReceivedMessage contains logic for creating a todo.
func (c Consumer) ReceivedMessage(message *sarama.ConsumerMessage) error {
	headers := make(map[string]string, len(message.Headers))
	for _, header := range message.Headers {
		headers[string(header.Key)] = string(header.Value)
	}

	var span opentracing.Span

	spanCtx, err := c.tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
	if err == nil {
		span = opentracing.GlobalTracer().StartSpan(spanName, opentracing.FollowsFrom(spanCtx))
	} else {
		span = opentracing.StartSpan(spanName)
		log.Printf("could not create span: %v", err)
	}

	defer span.Finish()

	var t todov1.CreateRequest
	if err := proto.Unmarshal(message.Value, &t); err != nil {
		return fmt.Errorf("could not deserialise todo: %v", err)
	}

	if err := c.repo.Create(
		opentracing.ContextWithSpan(context.Background(), span),
		&repository.Todo{
			Message: t.Message,
		},
	); err != nil {
		return fmt.Errorf("could not create todo, skipping message: %v", err)
	}

	return nil
}
