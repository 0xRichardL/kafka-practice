---
description: 'Memory bank for Kafka fraud detection project: components, topics, schemas, and modules. Use when working on any part of the fraud detection system.'
applyTo: '**'
---

# Kafka Fraud Detection Project Memory Bank

## Project Overview

This is a multi-service Kafka-based fraud detection system demonstrating real-world Kafka consumer patterns, state management, and failure handling. It intentionally avoids Kafka Streams and ksqlDB to focus on low-level consumer implementation.

## Services/Modules

- **producer/**: Golang service that generates and produces transaction events to Kafka
- **fraud/**: Golang service that consumes transactions, implements fraud detection logic, and produces alerts
- **alert/**: NestJS/Node.js service that consumes fraud alerts and provides REST API/WebSocket endpoints
- **fraud-kafka-streams/**: Java Kafka Streams application (alternative implementation)
- **schemas/**: Protobuf schema definitions and generated Go code for TransactionEvent and FraudAlertEvent

## Kafka Topics

- `transactions.raw`: Raw transaction events from producer
- `fraud.alerts`: Fraud detection alerts from fraud service
- `transactions.store`: Changelog topic for state store (from Kafka Streams)

## Message Schemas

- **TransactionEvent** (transaction.proto): Defines transaction data structure
- **FraudAlertEvent** (fraud_alert.proto): Defines fraud alert data structure
- Schemas use Protobuf with backward compatibility

## Infrastructure (docker-compose.yml)

- Kafka (KRaft mode, no Zookeeper)
- Schema Registry
- Kafka Connect (for sinks to PostgreSQL/Elasticsearch)
- Kafka UI (monitoring interface)

## Key Concepts

- Focus on consumption correctness, state management, failure modes, replayability
- Fraud detection rules: velocity, spike detection, country/MCC validation, impossible travel
- State stores: rolling averages, last location tracking
- Windowed aggregations for time-based rules

## Technologies

- Go (producer, fraud services)
- Node.js/NestJS (alert service)
- Java (Kafka Streams alternative)
- Protobuf (schemas)
- Docker Compose (infrastructure)
