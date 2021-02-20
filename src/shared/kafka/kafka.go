package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// Client wraps a sarama client.
type Client struct {
	saramaClient sarama.Client
}

// SyncProducer represents a sync producer.
type SyncProducer struct {
	producer sarama.SyncProducer
}

type ConsumerGroup struct {
	cGroup sarama.ConsumerGroup
}

// NewClient returns a new client. It retries forever until the connection is ready.
func NewClient(
	brokerAddresses []string,
	config *sarama.Config,
	waitFor time.Duration,
) (Client, error) {
	saramaClient, err := sarama.NewClient(brokerAddresses, config)
	if err != nil {
		log.Printf("kafka client not ready, retrying: %v", err)
		time.Sleep(waitFor)
		return NewClient(brokerAddresses, config, waitFor)
	}

	ctx, cancel := context.WithTimeout(context.Background(), waitFor)
	defer cancel()

	client := Client{saramaClient: saramaClient}

	if err := client.Ping(ctx); err != nil {
		log.Printf("kafka client connection not ready, retrying: %v", err)
		time.Sleep(waitFor)
		return NewClient(brokerAddresses, config, waitFor)
	}

	return client, nil
}

// NewSyncProducer returns a new sync producer.
func NewSyncProducer(client Client) (SyncProducer, error) {
	producer, err := sarama.NewSyncProducerFromClient(client.saramaClient)
	if err != nil {
		return SyncProducer{}, fmt.Errorf("could not create a new sync producer: %w", err)
	}
	return SyncProducer{producer: producer}, nil
}

// NewConsumerGroup returns a new consumer group.
func NewConsumerGroup(groupID string, client Client) (ConsumerGroup, error) {
	cGroup, err := sarama.NewConsumerGroupFromClient(groupID, client.saramaClient)
	if err != nil {
		return ConsumerGroup{}, fmt.Errorf("could not create a new consumer group: %w", err)
	}
	return ConsumerGroup{cGroup: cGroup}, nil
}

// Ping checks the status of the kafka connection.
func (c Client) Ping(ctx context.Context) error {
	if len(c.saramaClient.Brokers()) == 0 {
		return errors.New("could not ping kafka")
	}
	return nil
}

// SendMessage wraps the send message method.
func (sp SyncProducer) SendMessage(message *sarama.ProducerMessage) error {
	_, _, err := sp.producer.SendMessage(message)
	return err
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
