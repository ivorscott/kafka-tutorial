package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	bootstrap := getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka-cluster-kafka-bootstrap.kafka:9092")
	topic := getenv("TOPIC", "my-topic")

	log.Printf("Starting producer, bootstrap=%s topic=%s", bootstrap, topic)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	// Delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v", ev.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to %v", ev.TopicPartition)
				}
			}
		}
	}()

	// Ctrl+C handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-ticker.C:
			i++
			value := fmt.Sprintf("Hello from Go producer #%d at %s", i, time.Now().Format(time.RFC3339))
			err := p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(value),
			}, nil)
			if err != nil {
				log.Printf("Produce error: %v", err)
			}
		case <-sigCh:
			log.Println("Shutting down producer...")
			p.Flush(5000)
			return
		}
	}
}
