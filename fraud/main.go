package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Kafka
	kafkaService, err := NewKafkaService()
	if err != nil {
		log.Fatalf("Failed to create KafkaService: %s", err.Error())
	}
	defer kafkaService.Close()

	fraudEngine, err := NewFraudEngine(kafkaService)
	if err != nil {
		log.Fatalf("Failed to create FraudEngine: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()
	fraudEngine.StartIngestion(ctx)
}
