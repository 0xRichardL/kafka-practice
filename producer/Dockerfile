# --- Build Stage ---
FROM golang:1.25.4-bookworm AS builder

WORKDIR /workspace

# Install librdkafka dev headers
RUN apt-get update && apt-get install -y librdkafka-dev pkg-config

# Copy required modules
COPY producer/go.work.override go.work
COPY go.work.sum go.work.sum
COPY producer/ producer/
COPY schemas/ schemas/
COPY fraud/ fraud/

RUN go work sync

WORKDIR /workspace/producer

# CGO required
RUN go build -o /app/main .

# --- Final Stage ---
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y librdkafka1 ca-certificates && rm -rf /var/lib/apt/lists/*

RUN groupadd -r appgroup && useradd -r -g appgroup appuser

WORKDIR /home/appuser/app

COPY --from=builder /app/main .

RUN chown -R appuser:appgroup /home/appuser/app

USER appuser

EXPOSE 8080

CMD ["./main"]