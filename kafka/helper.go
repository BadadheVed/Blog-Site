package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/yourname/blog-kafka/models"
)

type NotificationPayload struct {
	ChannelID uuid.UUID               `json:"channel_id"`
	AuthorID  uuid.UUID               `json:"author_id"`
	BlogID    uuid.UUID               `json:"blog_id"`
	BlogTitle string                  `json:"blog_title"`
	Type      models.NotificationType `json:"type"`
}

// Writer is a package-level writer for a topic.
// It must be initialized by calling InitKafkaWriter(...) before PublishNotification.
var Writer *kafka.Writer

// InitKafkaWriter creates and stores a kafka.Writer for the given brokers/topic.
// Call this once at app startup.
func InitKafkaWriter(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	log.Printf("[kafka] writer initialized for topic=%s brokers=%v\n", topic, brokers)
}

// PublishNotification serializes payload and writes to Kafka. Returns an error if writer is nil or write fails.
func PublishNotification(payload NotificationPayload) error {
	if Writer == nil {
		return errors.New("kafka writer is not initialized: call kafka.InitKafkaWriter(...) first")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{Value: data}

	// Use a context rather than nil
	if err := Writer.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("[kafka] write failed: %v", err)
		return err
	}

	log.Printf("[kafka] published message for blog=%s channel=%s", payload.BlogID.String(), payload.ChannelID.String())
	return nil
}
