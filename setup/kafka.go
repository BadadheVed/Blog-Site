package setup

import (
	"log"

	"github.com/yourname/blog-kafka/kafka"
	"github.com/yourname/blog-kafka/notifications"
)

// StartKafkaSetup initializes the Kafka producer and consumers
func StartKafkaSetup(
	MyWorkerPool *notifications.WorkerPool,
	MyNotifService *notifications.NotificationService,
) (*kafka.Producer, error) {
	// Set the global references in the kafka package
	kafka.SetWorkerPoolAndService(MyWorkerPool, MyNotifService)

	brokers := []string{"localhost:9092"}
	hotTopic := "hot-notifications"
	coldTopic := "cold-notifications"
	groupID := "notification-group"

	// Initialize the Kafka producer
	producer := kafka.NewProducer(brokers)

	// Start consumers
	go kafka.StartConsumer(brokers, groupID, hotTopic)
	go kafka.StartConsumer(brokers, groupID, coldTopic)

	log.Println("[kafka] producer initialized & hot/cold consumers started")

	return producer, nil
}
