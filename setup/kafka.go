package setup

import (
	"log"

	"github.com/yourname/blog-kafka/function"
	"github.com/yourname/blog-kafka/kafka"
	"github.com/yourname/blog-kafka/notifications"
)

func StartKafkaSetup(MyWorkerPool *notifications.WorkerPool, MyNotifSvc *notifications.NotificationService) {
	function.SetWorkerPool(MyWorkerPool)

	brokers := []string{"localhost:9092"}
	topic := "blog-events"
	groupID := "notification-group"

	kafka.InitKafkaWriter(brokers, topic) // optional if producer needed
	go kafka.StartBlogEventConsumer(MyNotifSvc, brokers, groupID, topic)

	log.Println("[kafka] blog event consumer started")
}
