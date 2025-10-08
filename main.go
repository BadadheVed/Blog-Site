package main

import (
	"fmt"
	"time"

	"github.com/yourname/blog-kafka/notifications"
	"github.com/yourname/blog-kafka/workers"
)

var MyWorkerPool *workers.WorkerPool

func InitWorkerPool() {
	handler := &notifications.NotificationService{}
	MyWorkerPool = workers.NewWorkerPool(5, 100, 30*time.Second, handler)
	fmt.Println("Go Workers Started")
}

// func main() {
// 	config.DBConnect()
// 	InitWorkerPool()
// 	if MyWorkerPool == nil {
// 		log.Fatal("Worker pool initialization failed!")
// 	}
// 	function.SetWorkerPool(MyWorkerPool)

// 	r := routes.SetupRouter()
// 	port := os.Getenv("PORT")

// 	if port == "" {
// 		port = "8080"
// 	}

// 	log.Printf("Server running on port %s", port)
// 	r.Run(":" + port)

// }

// func main() {
// 	brokers := []string{"localhost:9092"}

// 	config := sarama.NewConfig()
// 	config.Producer.Return.Successes = true
// 	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
// 	prod, err := sarama.NewSyncProducer(brokers, config)
// 	if err != nil {
// 		log.Fatalf("Failed to start producer: %v", err)
// 	}
// 	defer prod.Close()

// 	for i := 0; i < 10; i++ {
// 		msg := &sarama.ProducerMessage{
// 			Topic: "test-topic",
// 			Value: sarama.StringEncoder("Hello Kafka! #" + string(i)),
// 		}

// 		partition, offset, err := prod.SendMessage(msg)
// 		if err != nil {
// 			log.Fatalf("Failed to send message: %v", err)
// 		}
// 		log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
// 	}

// }
