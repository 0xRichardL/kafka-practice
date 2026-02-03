package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

var (
	defaultConfig = kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
		"group.id":          os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
		"auto.offset.reset": "earliest",
		"security.protocol": "PLAINTEXT",
	}
)

type KafkaService struct {
	Producer     *kafka.Producer
	Consumer     *kafka.Consumer
	Serializer   *protobuf.Serializer
	Deserializer *protobuf.Deserializer
}

func NewKafkaService() (*KafkaService, error) {
	// Schema Registry
	schemaRegistryClient, err := schemaregistry.NewClient(
		schemaregistry.NewConfig(os.Getenv("SCHEMA_REGISTRY_URL")),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Kafka schema registry: %w", err)
	}

	// Kafka Serializer & Deserializer
	serializer, err := protobuf.NewSerializer(
		schemaRegistryClient,
		serde.ValueSerde,
		protobuf.NewSerializerConfig(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create serializer: %w", err)
	}
	deserializer, err := protobuf.NewDeserializer(
		schemaRegistryClient,
		serde.ValueSerde,
		protobuf.NewDeserializerConfig(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create deserializer: %w", err)
	}

	// Kafka Producer & Consumer
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

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
		"group.id":          os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
		"auto.offset.reset": "earliest",
		"security.protocol": "PLAINTEXT",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &KafkaService{
		producer,
		consumer,
		serializer,
		deserializer,
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

func (k *KafkaService) CreateConsumer(config *kafka.ConfigMap) (*kafka.Consumer, error) {
	for k, v := range defaultConfig {
		if _, exists := (*config)[k]; !exists {
			(*config)[k] = v
		}
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	return consumer, nil
}

func (k *KafkaService) Close() {
	k.Producer.Flush(15 * 1000) // 15s
	k.Producer.Close()
	if err := k.Consumer.Close(); err != nil {
		fmt.Printf("Failed to close Kafka consumer: %v\n", err)
	}
}
