package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

type TransactionService struct {
	kafkaProducer *kafka.Producer
	serializer    *protobuf.Serializer
	isOn          bool
}

func NewTxProducer(brokers string, serializer *protobuf.Serializer) (*TransactionService, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"enable.idempotence": true,
		"compression.codec":  "snappy",
		"security.protocol":  "SASL_SSL",
		"sasl.mechanisms":    "PLAIN",
		"acks":               "all",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	go func() {
		for e := range producer.Events() {
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
	}()

	return &TransactionService{
		kafkaProducer: producer,
		serializer:    serializer,
	}, nil
}

func (s *TransactionService) Close() {
	if s.kafkaProducer == nil {
		return
	}
	s.kafkaProducer.Flush(15 * 1000) // 15s
	s.kafkaProducer.Close()
}

func (s *TransactionService) Transfer() {
	// TODO:
	// 1. gen Proto go code.
	// 2. serialize by registry.
	// 3. Send Kafka.
}

func (s *TransactionService) GenerateTransfers(num int) {

}
