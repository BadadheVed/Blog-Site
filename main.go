package main

import (
	"log"
	"os"
	"time"

	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/function"
	"github.com/yourname/blog-kafka/notifications"
	"github.com/yourname/blog-kafka/routes"
	"github.com/yourname/blog-kafka/setup"
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

	producer, err := setup.StartKafkaSetup(MyWorkerPool, MyNotifSvc)
	if err != nil {
		log.Fatalf("[kafka] Failed to setup Kafka: %v", err)
	}
	function.InitKafkaProducer(producer)
	log.Println("[init] Kafka producer initialized in function package")

	r := routes.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[server] Listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[server] Failed to run: %v", err)
	}

}
