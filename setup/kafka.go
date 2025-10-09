package setup

import (
	"log"
	"os"
	"strings"

	"github.com/yourname/blog-kafka/kafka"
	"github.com/yourname/blog-kafka/notifications"
)

func StartKafkaSetup(
	MyWorkerPool *notifications.WorkerPool,
	MyNotifService *notifications.NotificationService,
) (*kafka.Producer, error) {

	kafka.SetWorkerPoolAndService(MyWorkerPool, MyNotifService)

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092"
	}
	brokers := strings.Split(brokersEnv, ",")

	hotTopic := "hot-notifications"
	coldTopic := "cold-notifications"
	groupID := "notification-group"

	producer := kafka.NewProducer(brokers)

	go kafka.StartConsumer(brokers, groupID, hotTopic)
	go kafka.StartConsumer(brokers, groupID, coldTopic)

	log.Printf("[kafka] producer initialized & hot/cold consumers started on brokers=%v\n", brokers)

	return producer, nil
}
