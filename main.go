package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/function"
	"github.com/yourname/blog-kafka/notifications"
	"github.com/yourname/blog-kafka/routes"
	"github.com/yourname/blog-kafka/workers"
)

var MyWorkerPool *workers.WorkerPool

func InitWorkerPool() {
	handler := &notifications.NotificationService{}
	MyWorkerPool = workers.NewWorkerPool(5, 100, 30*time.Second, handler)
	fmt.Println("Go Workers Started")
}

func main() {
	config.DBConnect()
	InitWorkerPool()
	if MyWorkerPool == nil {
		log.Fatal("Worker pool initialization failed!")
	}
	function.SetWorkerPool(MyWorkerPool)

	r := routes.SetupRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)

}
