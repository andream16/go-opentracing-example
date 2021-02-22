package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

// ConsumerGroup wraps a sarama consumer group.
type ConsumerGroup struct {
	cGroup sarama.ConsumerGroup
}

// NewConsumerGroup returns a new consumer group.
func NewConsumerGroup(groupID string, client Client) (ConsumerGroup, error) {
	cGroup, err := sarama.NewConsumerGroupFromClient(groupID, client.saramaClient)
	if err != nil {
		return ConsumerGroup{}, fmt.Errorf("could not create a new consumer group: %w", err)
	}
	return ConsumerGroup{cGroup: cGroup}, nil
}

// Consume consumes from the given topics and calls the passed handler.
func (c ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return c.cGroup.Consume(ctx, topics, handler)
}

// Errors returns the returned errors from the consumer group.
func (c ConsumerGroup) Errors() <-chan error {
	return c.cGroup.Errors()
}

// Close closes the consumer group.
func (c ConsumerGroup) Close() error {
	return c.cGroup.Close()
}
