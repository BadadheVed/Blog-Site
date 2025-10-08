package main

import (
	"log"
	"os"
	"time"

	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/function"
	"github.com/yourname/blog-kafka/kafka"
	"github.com/yourname/blog-kafka/notifications"
	"github.com/yourname/blog-kafka/routes"
)

var (
	MyWorkerPool *notifications.WorkerPool
	MyNotifSvc   *notifications.NotificationService
)

func InitWorkerPool() {

	handler := &notifications.NotificationService{}

	MyWorkerPool = notifications.NewWorkerPool(100, 30*time.Second, handler)

	handler.WorkerPool = MyWorkerPool
	MyNotifSvc = handler

	log.Println("[init] worker pool + notification service initialized")
}

func main() {

	config.DBConnect()

	InitWorkerPool()

	function.SetWorkerPool(MyWorkerPool)

	brokers := []string{"localhost:9092"}
	groupID := "notification-group"
	topic := "blog-events"

	go kafka.StartBlogEventConsumer(MyNotifSvc, brokers, groupID, topic)
	log.Println("[kafka] blog event consumer started")

	r := routes.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[server] listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[server] failed to run: %v", err)
	}
}
