package main

import (
	"context"
	"fmt"
	"log"

	"github.com/0xRichardL/kafka-practice/schemas"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	RAW_TX_TOPIC string = "transactions.raw"
)

type FraudEngine struct {
	kafkaService  *KafkaService
	rawTxConsumer *kafka.Consumer
}

func NewFraudEngine(kafkaService *KafkaService) (*FraudEngine, error) {
	rawTxConsumer, err := kafkaService.CreateConsumer(&kafka.ConfigMap{
		"group.id":           "fraud-engine",
		"enable.auto.commit": false,
		"session.timeout.ms": 10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create raw transaction consumer: %w", err)
	}
	return &FraudEngine{
		kafkaService:  kafkaService,
		rawTxConsumer: rawTxConsumer,
	}, nil
}

func (f *FraudEngine) StartIngestion(ctx context.Context) error {
	if err := f.rawTxConsumer.SubscribeTopics([]string{RAW_TX_TOPIC}, nil); err != nil {
		return fmt.Errorf("failed to subscribe to raw transactions topic: %w", err)
	}
	defer f.rawTxConsumer.Close()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			event := f.rawTxConsumer.Poll(500)
			if event == nil {
				continue
			}
			switch e := event.(type) {
			case *kafka.Message:
				data, err := f.kafkaService.Deserializer.Deserialize(RAW_TX_TOPIC, e.Value)
				if err != nil {
					log.Printf("failed to deserialize message: %v\n", err)
					continue
				}
				tx, ok := data.(*schemas.Transaction)
				if !ok {
					log.Printf("unexpected data type: %T\n", data)
					continue
				}
				if err := f.ingestTransaction(tx); err != nil {
					log.Printf("failed to ingest transaction %s: %v\n", tx.TransactionId, err)
					continue
				}
				if _, err := f.rawTxConsumer.CommitMessage(e); err != nil {
					log.Printf("failed to commit message: %v\n", err)
				}

			case kafka.Error:
				log.Printf("kafka error: %v\n", e)

			case kafka.AssignedPartitions:
				log.Printf("assigned: %v\n", e.Partitions)
				f.rawTxConsumer.Assign(e.Partitions)

			case kafka.RevokedPartitions:
				log.Printf("revoked: %v\n", e.Partitions)
				f.rawTxConsumer.Unassign()
			}

		}
		return nil
	}
}

func (f *FraudEngine) ingestTransaction(tx *schemas.Transaction) error {
	return nil
}
