package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
)

type Consumer struct {
	tracer opentracing.Tracer
}

// NewConsumer returns a new consumer.
func NewConsumer() Consumer {
	return Consumer{tracer: opentracing.GlobalTracer()}
}

func (c Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		span := c.tracer.StartSpan("kafka-consumer")
		defer span.Finish()

		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
