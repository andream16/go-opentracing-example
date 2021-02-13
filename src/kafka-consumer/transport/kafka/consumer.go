package kafka

import (
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"

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

		h := make(http.Header, len(message.Headers))
		for _, header := range message.Headers {
			h[string(header.Key)] = []string{string(header.Value)}
		}

		gt := opentracing.GlobalTracer()

		spanCtx, _ := gt.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(h))
		span := gt.StartSpan("consumer", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		session.MarkMessage(message, "")
	}

	return nil
}
