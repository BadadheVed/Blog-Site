package main

import (
	"fmt"

	"github.com/yourname/blog-kafka/config"
)

func main() {
	config.DBConnect()
	fmt.Println("Database connected sucessfully")
}
