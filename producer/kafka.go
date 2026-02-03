package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

type KafkaService struct {
	Producer   *kafka.Producer
	Serializer *protobuf.Serializer
}

func NewKafkaService() (*KafkaService, error) {
	/// Schema Registry
	schemaRegistryClient, err := schemaregistry.NewClient(
		schemaregistry.NewConfig(os.Getenv("SCHEMA_REGISTRY_URL")),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Kafka schema registry: %w", err)
	}

	/// Kafka Serializer
	serializer, err := protobuf.NewSerializer(
		schemaRegistryClient,
		serde.ValueSerde,
		protobuf.NewSerializerConfig(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create serializer: %w", err)
	}

	/// Kafka Producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_BROKERS"),
		"enable.idempotence": true,
		"compression.codec":  "snappy",
		"security.protocol":  "PLAINTEXT",
		"acks":               "all",
		// Retry configurations
		// Ordering is preserved when enable.idempotence = true
		"retries":              5,
		"retry.backoff.ms":     200,
		"retry.backoff.max.ms": 10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &KafkaService{
		producer,
		serializer,
	}, nil
}

// BackgroundProduce is for handle async/background events
// which not being wait by the emitter (no delivery channel passed).
func (k *KafkaService) BackgroundProduce() {
	for e := range k.Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			}
		}
	}
}

func (k *KafkaService) Close() {
	k.Producer.Flush(15 * 1000) // 15s
	k.Producer.Close()
}
