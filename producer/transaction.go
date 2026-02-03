package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/0xRichardL/kafka-practice/schemas"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type TransactionService struct {
	kafkaService *KafkaService
	kafkaTopic   string
	control      bool
	mu           *sync.Mutex
	currentTxID  int
}

func NewTransactionService(kafkaService *KafkaService) *TransactionService {
	return &TransactionService{
		kafkaService: kafkaService,
		kafkaTopic:   "transactions.raw",
		mu:           &sync.Mutex{},
	}
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
	payload, err := s.kafkaService.Serializer.Serialize(s.kafkaTopic, &tx)
	if err != nil {
		return fmt.Errorf("failed to serialize: %w", err)
	}

	err = s.kafkaService.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.kafkaTopic,
			Partition: kafka.PartitionAny, // let Kafka decide the partition by key
		},
		Key:   []byte(tx.UserId), // partition by UserId
		Value: payload,
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
