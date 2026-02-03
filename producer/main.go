package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Kafka
	kafkaService, err := NewKafkaService()
	if err != nil {
		log.Fatalf("Failed to create KafkaService: %s", err.Error())
	}
	go kafkaService.BackgroundProduce()
	defer kafkaService.Close()
	// Services
	transactionService := NewTransactionService(kafkaService)
	// REST
	router := gin.Default()
	controller := NewController(transactionService)
	controller.RegisterRoutes(router)
	log.Println("Starting Fraud server on port 5001")
	if err := router.Run(":5001"); err != nil {
		log.Fatalln("Unable to start Fraud server", err)
	}
}
