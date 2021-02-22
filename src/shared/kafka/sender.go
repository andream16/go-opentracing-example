package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// Sender describe the send contract.
type Sender interface {
	SendMessage(message *sarama.ProducerMessage) error
}

// SyncProducer represents a sync producer.
type SyncProducer struct {
	producer sarama.SyncProducer
}

// NewSyncProducer returns a new sync producer.
func NewSyncProducer(client Client) (SyncProducer, error) {
	producer, err := sarama.NewSyncProducerFromClient(client.saramaClient)
	if err != nil {
		return SyncProducer{}, fmt.Errorf("could not create a new sync producer: %w", err)
	}
	return SyncProducer{producer: producer}, nil
}

// SendMessage wraps the send message method.
func (sp SyncProducer) SendMessage(message *sarama.ProducerMessage) error {
	_, _, err := sp.producer.SendMessage(message)
	return err
}
