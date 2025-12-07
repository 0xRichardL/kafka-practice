package main

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"github.com/gin-gonic/gin"
)

func main() {
	// Kafka
	schemaRegistryClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(os.Getenv("SCHEMA_REGISTRY_URL")))
	if err != nil {
		log.Fatalf("Failed to connect to Kafka schema registry: %s", err.Error())
	}
	serializer, err := protobuf.NewSerializer(schemaRegistryClient, serde.ValueSerde, protobuf.NewSerializerConfig())
	if err != nil {
		log.Fatalf("Failed to create serializer: %s", err.Error())
	}
	txProducer, err := NewTxProducer(os.Getenv("KAFKA_BROKERS"), serializer)
	if err != nil {
		log.Fatalf("Failed to create TxProducer: %s", err.Error())
	}
	// REST
	router := gin.Default()
	controller := NewController(txProducer)
	controller.RegisterRoutes(router)
	log.Println("Starting Fraud server on port 5001")
	if err := router.Run(":5001"); err != nil {
		log.Fatalln("Unable to start Fraud server", err)
	}
}
