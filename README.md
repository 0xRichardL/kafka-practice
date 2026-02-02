# Kafka Fraud Detection Project — Implementation Checklist

## 1. 📦 Environment & Infra Setup

- [x] Install Docker & Docker Compose
- [x] Install Go
- [x] Install Node.js
- [x] Create project root folder
- [x] Add docker-compose.yml containing:
  - [x] Kafka (KRaft mode)
  - [ ] Zookeeper (if not using KRaft)
  - [x] Schema Registry
  - [x] Kafka Connect
  - [x] Kafka UI / Conduktor
  - [ ] Postgres / Elasticsearch
- [x] Confirm topics can be created
- [x] Check Kafka UI connectivity
- [x] Verify Schema Registry is reachable (/subjects endpoint)

## 2. 📚 Define Message Schemas

- [x] Create `/schemas` folder
- [x] Write Avro/Protobuf schema for TransactionEvent
- [x] Write schema for FraudAlertEvent
- [x] Enable schema compatibility rules (BACKWARD or FULL)
- [ ] Add schema files to version control

## 3. 🏭 Build the Producer Service (Go)

- [x] Create `/producer` folder
- [x] Add config env: broker URL, schema-registry URL
- [x] Integrate Avro/Protobuf serializer
- [x] Auto-register schema on startup
- [x] Produce random or API-driven transaction events
- [x] Add retry + backoff for production failures
- [ ] Add structured logs
- [x] Add Dockerfile
- [x] Validate messages appear in Kafka

## 4. 🔍 Fraud Detection Service (Go)

- [x] Create `/fraud` folder
- [ ] Implement Kafka consumer:
  - [ ] Consume from `transactions` topic
  - [ ] Handle schema deserialization (Protobuf)
  - [ ] Implement fraud detection logic
  - [ ] Maintain state/metrics in memory or external store
  - [ ] Produce alerts to `fraud-alerts` topic
- [ ] Add structured logging
- [ ] Add metrics (Prometheus)
- [ ] Add Dockerfile
- [ ] Validate alerts are produced correctly

## What This Project Intentionally Does NOT Use

- Kafka Streams (Java)
- ksqlDB
- Exactly-once semantics at broker level

This is a deliberate design choice to demonstrate
how Kafka consumers work under real-world constraints.

## Why This Project Matters

Most Kafka tutorials stop at producing messages.

This project focuses on:

- consumption correctness
- state
- failure modes
- replayability

Which is where senior-level Kafka understanding lives.

## Future Improvements

- Persistent state backend
- Metrics & tracing
- Dead-letter topic
- Dynamic rule loading
- ksqlDB integration for simple rules
