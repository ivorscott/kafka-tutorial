# Build stage
FROM golang:1.22-bullseye AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# librdkafka dev libs for build
RUN apt-get update && apt-get install -y --no-install-recommends \
    librdkafka-dev pkg-config && \
    rm -rf /var/lib/apt/lists/*

RUN CGO_ENABLED=1 go build -o producer ./...

# Runtime stage
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates librdkafka1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/producer /app/producer

ENV KAFKA_BOOTSTRAP_SERVERS=kafka-cluster-kafka-bootstrap.kafka:9092
ENV TOPIC=my-topic

CMD ["/app/producer"]