package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/0xRichardL/kafka-practice/schemas"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

type TransactionService struct {
	kafkaProducer *kafka.Producer
	kafkaTopic    string
	serializer    *protobuf.Serializer
	control       bool
	mu            *sync.Mutex
	currentTxID   int
}

func NewTxProducer(brokers string, serializer *protobuf.Serializer) (*TransactionService, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"enable.idempotence": true,
		"compression.codec":  "snappy",
		"security.protocol":  "PLAINTEXT",
		"acks":               "all",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// This is for handle async/background events which not being wait by the emitter (no delivery channel passed).
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
		kafkaTopic:    "transactions.raw",
		serializer:    serializer,
		mu:            &sync.Mutex{},
	}, nil
}

func (s *TransactionService) Close() {
	if s.kafkaProducer == nil {
		return
	}
	s.kafkaProducer.Flush(15 * 1000) // 15s
	s.kafkaProducer.Close()
}

func (s *TransactionService) genTxID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTxID++
	return fmt.Sprint(s.currentTxID)
}

func (s *TransactionService) Transfer() error {
	tx := schemas.Transaction{
		TransactionId: s.genTxID(),
		UserId:        fmt.Sprint(rand.Intn(100)),
		Amount:        rand.Float64() + float64(rand.Intn(1000)),
		Merchant:      fmt.Sprint(rand.Intn(100)),
		Category:      fmt.Sprint(rand.Intn(9)),
		Timestamp:     time.Now().Unix(),
		Status:        "initialized",
	}
	payload, err := s.serializer.Serialize(s.kafkaTopic, &tx)
	if err != nil {
		return fmt.Errorf("failed to serialize: %w", err)
	}

	err = s.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.kafkaTopic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce: %w", err)
	}
	return nil
}

func (s *TransactionService) GenerateTransfers(control bool) {
	s.control = control
	go func() {
		for i := 0; s.control && i < 1000000; i++ {
			if err := s.Transfer(); err != nil {
				fmt.Printf("failed to transfer: %s", err.Error())
				return
			}
		}
	}()
}
