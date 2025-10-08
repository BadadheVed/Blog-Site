package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/yourname/blog-kafka/models"
	"github.com/yourname/blog-kafka/notifications"
)

type NotificationEvent struct {
	ChannelID uuid.UUID               `json:"channel_id"`
	AuthorID  uuid.UUID               `json:"author_id"`
	BlogID    uuid.UUID               `json:"blog_id"`
	BlogTitle string                  `json:"blog_title"`
	Type      models.NotificationType `json:"type"`
}

type Consumer struct {
	Reader     *kafka.Reader
	NotifSvc   *notifications.NotificationService
	TopicLabel string
}

func NewConsumer(brokers []string, topic, groupID, label string, notifSvc *notifications.NotificationService) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		Reader:     r,
		NotifSvc:   notifSvc,
		TopicLabel: label,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka consumer for topic [%s]", c.TopicLabel)
	for {
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			log.Println("FetchMessage error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var event NotificationEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Invalid event: %v", err)
			continue
		}

		log.Printf("[%s] Received notification event: %+v", c.TopicLabel, event)

		err = c.NotifSvc.CreateNotificationForChannelMembers(
			event.ChannelID,
			event.AuthorID,
			event.BlogID,
			event.BlogTitle,
			event.Type,
		)
		if err != nil {
			log.Printf("Error creating notifications: %v", err)
		}

		c.Reader.CommitMessages(ctx, m)
	}
}
