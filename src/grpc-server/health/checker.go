package health

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"

	"github.com/andream16/go-opentracing-example/src/shared/health"
)

// KafkaChecker represents a kafka checker.
type KafkaChecker struct {
	client sarama.Client
}

type UnhealthyConnectionError struct {
	err error
}

func (u UnhealthyConnectionError) Error() string {
	return u.err.Error()
}

func (u UnhealthyConnectionError) IsRetriable() bool {
	return true
}

// NewKafkaChecker returns a new kafka checker.
func NewKafkaChecker(client sarama.Client) health.Checker {
	return KafkaChecker{client: client}
}

// Check checks the underlying kafka connection.
func (k KafkaChecker) Check(ctx context.Context) error {
	if len(k.client.Brokers()) == 0 {
		return errors.New("no brokers found")
	}
	return nil
}

func (k KafkaChecker) Name() string {
	return "kafka"
}
