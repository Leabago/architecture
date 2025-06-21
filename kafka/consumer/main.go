package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "my-topic"
	// make a new reader that consumes from topic, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   "my-groupID",
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	// r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		time.Sleep(5 * time.Second)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
