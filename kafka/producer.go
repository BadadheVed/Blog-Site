package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/models"
)

type NotificationPayload struct {
	ChannelID uuid.UUID               `json:"channel_id"`
	AuthorID  uuid.UUID               `json:"author_id"`
	BlogID    uuid.UUID               `json:"blog_id"`
	BlogTitle string                  `json:"blog_title"`
	Type      models.NotificationType `json:"type"`
}

type Producer struct {
	SyncProducer sarama.SyncProducer
}

func NewProducer(brokers []string) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	prod, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("[kafka] failed to start Sarama producer: %v", err)
	}

	log.Printf("[kafka] producer initialized for brokers=%v", brokers)
	return &Producer{SyncProducer: prod}
}

func (p *Producer) Send(topic, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("[kafka] failed to send message to %s: %v", topic, err)
		return
	}

	log.Printf("[kafka] message sent to %s partition=%d offset=%d", topic, partition, offset)
}

func (p *Producer) PublishNotification(payload NotificationPayload, hot bool) {
	topic := "cold-notifications"
	if hot {
		topic = "hot-notifications"
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[kafka] failed to marshal payload: %v", err)
		return
	}

	p.Send(topic, string(data))
}
