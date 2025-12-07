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
- [ ] Confirm topics can be created
- [ ] Check Kafka UI connectivity
- [ ] Verify Schema Registry is reachable (/subjects endpoint)

## 2. 📚 Define Message Schemas

- [x] Create `/schemas` folder
- [x] Write Avro/Protobuf schema for TransactionEvent
- [x] Write schema for FraudAlertEvent
- [x] Enable schema compatibility rules (BACKWARD or FULL)
- [ ] Add schema files to version control

## 3. 🏭 Build the Producer Service (Go)

- [x] Create `/producer` folder
- [x] Add config env: broker URL, schema-registry URL
- [ ] Integrate Avro/Protobuf serializer
- [ ] Auto-register schema on startup
- [ ] Produce random or API-driven transaction events
- [ ] Add retry + backoff for production failures
- [ ] Add structured logs
- [ ] Add Dockerfile
- [ ] Validate messages appear in Kafka

## 4. Fraud Detection Processor (Kafka Streams)

- [ ] Create /services/fraud-detector folder
- [ ] Implement Streams topology:
  - [ ] Ingest `transactions` topic
  - [ ] Use state store (RocksDB)
  - [ ] Compute running metrics
  - [ ] Detect fraud patterns
  - [ ] Emit to `fraud-alerts` topic
- [ ] Handle schema deserialization
- [ ] Add metrics (Prometheus)
- [ ] Add Dockerfile
- [ ] Validate alerts are produced
