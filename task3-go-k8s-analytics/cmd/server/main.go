package main

import (
	"log"

	"task3-go-k8s-analytics/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
