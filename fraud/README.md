# Kafka Fraud Detection Consumer (Go)

A production-style Kafka consumer written in Go that demonstrates **real Kafka fundamentals**:
stateful consumption, idempotent processing, replayability, and failure handling.

This service consumes transaction events, applies fraud detection rules, and emits alerts —
focusing on **consumer-side complexity**, where Kafka skills actually matter.

## Architecture Overview

```
Kafka (transactions.raw)
        ↓
Fraud Consumer (Go)
  - Consumer Group
  - Manual Offsets
  - Stateful Processing
  - Idempotency
        ↓
Kafka (fraud.alerts)
```

## Core Responsibilities

- Consume transaction events from Kafka
- Maintain per-entity state (rolling windows, counters)
- Detect fraud patterns
- Emit fraud alerts
- Survive crashes, rebalances, and reprocessing

## Key Kafka Concepts Demonstrated

| Concept                | How It’s Implemented                      |
| ---------------------- | ----------------------------------------- |
| Consumer Groups        | Shared group for horizontal scaling       |
| At-least-once Delivery | Manual offset commits                     |
| Idempotency            | Deduplication by transaction ID           |
| Ordering               | Keyed consumption per entity              |
| Stateful Processing    | In-memory + persistent state              |
| Rebalancing            | Graceful partition revocation handling    |
| Replay                 | Offset reset & deterministic reprocessing |
| Backpressure           | Controlled polling & processing           |

## Delivery Guarantees

- **Kafka guarantee:** At-least-once
- **Application guarantee:** Effectively-once (via idempotent logic)
- **External systems:** Idempotent writes required

Exactly-once semantics are not available in Go clients; this service demonstrates
how real systems achieve correctness without Kafka Streams.

## Fraud Detection Logic (Example)

Implemented rules may include:

- Velocity spikes (N transactions per time window)
- Amount anomalies
- Impossible travel patterns
- Repeated merchant/category usage

Rules are deterministic, enabling safe replay.

## State Management

State is maintained per key (e.g. user_id):

- Rolling time windows
- Counters and aggregates
- TTL-based eviction
- Optional persistence (Redis / Badger / Pebble)

State is rebuilt automatically during replay.

## Failure & Recovery Handling

- Graceful shutdown (SIGTERM)
- Manual offset commit after successful processing
- Safe handling of duplicate messages
- Automatic recovery via Kafka replay

## Local Development

### Prerequisites

- Kafka cluster
- Schema Registry (Protobuf)
- Go 1.21+

### Run

go run cmd/consumer/main.go

## Configuration

Key settings:

- KAFKA_BROKERS
- KAFKA_GROUP_ID
- KAFKA_INPUT_TOPIC
- KAFKA_OUTPUT_TOPIC
- STATE_BACKEND
