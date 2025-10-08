package kafka

import (
	"encoding/json"

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

var Writer *kafka.Writer

func InitKafkaWriter(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
}

func PublishNotification(payload NotificationPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	msg := kafka.Message{Value: data}
	return Writer.WriteMessages(nil, msg)
}
