package main

import (
	"log"
)

func main() {
	application, err := NewApp("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
