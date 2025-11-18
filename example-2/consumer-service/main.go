package main

import (
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
	group := getenv("GROUP_ID", "go-consumer-group")

	log.Printf("Starting consumer, bootstrap=%s topic=%s group=%s", bootstrap, topic, group)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for {
		select {
		case <-sigCh:
			log.Println("Shutting down consumer...")
			return
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("Received message at %s: %s",
					time.Now().Format(time.RFC3339), string(e.Value))
			case kafka.Error:
				log.Printf("Kafka error: %v", e)
			default:
				// ignore other events
			}
		}
	}
}
