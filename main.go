package main

import (
	"log"
	"os"

	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/routes"
)

func main() {
	config.DBConnect()
	r := routes.SetupRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port", port)
	r.Run(":" + port)

}
