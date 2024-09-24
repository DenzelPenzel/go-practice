package main

import (
	"testing"
	"time"
)

func TestPublishSubscribe(t *testing.T) {
	broker := NewBroker()
	defer broker.Close()

	sub, err := broker.Subscribe("test", 1)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	received := make(chan string)

	sub.Handler = func(msg Message) {
		if msg.Topic != "test" {
			t.Errorf("Expected topic 'test', got '%s'", msg.Topic)
		}
		payload, ok := msg.Payload.(string)
		if !ok {
			t.Errorf("Expected payload string, got %T", msg.Payload)
		}
		received <- payload
	}

	go func() {
		err := broker.Publish("test", "Hello, World!")
		if err != nil {
			t.Errorf("Failed to publish: %v", err)
		}
	}()

	select {
	case msg := <-received:
		if msg != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got '%s'", msg)
		}
	case <-time.After(1 * time.Second):
		t.Error("Did not receive message in time")
	}
}
