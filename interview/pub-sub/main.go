package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	broker := NewBroker()
	defer broker.Close()

	// Create a publisher
	publisher := NewPublisher(broker)

	// Subscribe to topics
	subscribeToTopic(broker, "news", "Subscriber 1")
	subscribeToTopic(broker, "sports", "Subscriber 2")
	subscribeToTopic(broker, "news", "Subscriber 3")

	// Handle graceful shutdown on interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Simulate publishing messages
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		count := 1
		for range ticker.C {
			publisher.Publish("news", fmt.Sprintf("Breaking News %d", count))
			publisher.Publish("sports", fmt.Sprintf("Sports Update %d", count))
			count++
			if count > 10 {
				return
			}
		}

	}()

	// Wait for interrupt signal
	<-sig
	log.Println("Interrupt signal received. Shutting down...")
}

// subscribeToTopic subscribes to a topic with a custom handler
func subscribeToTopic(broker *Broker, topic, name string) {
	// Buffer size of 10, don't blocking the broker
	sub, err := broker.Subscribe(topic, 10)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic '%s': %v", topic, err)
	}

	sub.Handler = func(msg Message) {
		log.Printf("[%s] Received on '%s': %v\n", name, msg.Topic, msg.Payload)
	}

	log.Printf("%s subscribed to topic: %s\n", name, topic)
}
