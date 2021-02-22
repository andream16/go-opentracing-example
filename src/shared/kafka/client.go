package kafka

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// Client wraps a sarama client.
type Client struct {
	saramaClient sarama.Client
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

// Ping checks the status of the kafka connection.
func (c Client) Ping(ctx context.Context) error {
	if len(c.saramaClient.Brokers()) == 0 {
		return errors.New("could not ping kafka")
	}
	return nil
}
