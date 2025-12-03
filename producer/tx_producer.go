package main

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type TxProducer struct {
	Writer *kafka.Writer
}

func (k *TxProducer) Close() {
	k.Writer.Close()
}

func NewTxProducer(brokers []string) *TxProducer {
	return &TxProducer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:          brokers,
			Topic:            "transactions",
			Balancer:         &kafka.LeastBytes{},
			RequiredAcks:     int(kafka.RequireAll),
			Async:            false,
			BatchTimeout:     10 * time.Millisecond,
			BatchSize:        100,
			BatchBytes:       1e6, // 1MB
			CompressionCodec: kafka.Snappy.Codec(),
			Logger:           kafka.LoggerFunc(log.Printf),
			ErrorLogger:      kafka.LoggerFunc(log.Printf),
		}),
	}
}
