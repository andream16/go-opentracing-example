package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
)

type Consumer struct{}

// NewConsumer returns a new consumer.
func NewConsumer() Consumer {
	return Consumer{}
}

func (c Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf(
			"Message claimed: value = %s, timestamp = %v, topic = %s\n",
			string(message.Value),
			message.Timestamp,
			message.Topic,
		)

		headers := make(map[string]string, len(message.Headers))
		for _, header := range message.Headers {
			headers[string(header.Key)] = string(header.Value)
		}

		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
		span := opentracing.GlobalTracer().StartSpan("consumer", opentracing.FollowsFrom(spanCtx))

		session.MarkMessage(message, "")

		span.Finish()
	}

	return nil
}
