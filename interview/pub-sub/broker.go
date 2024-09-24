package main

import (
	"context"
	"errors"
	"log"
	"sync"
)

// Broker manages subscriptions and publishing
type Broker struct {
	subscribers map[string]map[*Subscriber]struct{}
	lock        sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewBroker() *Broker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broker{
		subscribers: make(map[string]map[*Subscriber]struct{}),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Subscribe adds a subscriber to a topic.
// Returns the subscriber or an error.
func (b *Broker) Subscribe(topic string, buffer int) (*Subscriber, error) {
	if topic == "" {
		return nil, errors.New("topic cannot be empty")
	}

	sub := NewSubscriber(topic, buffer)

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make(map[*Subscriber]struct{})
	}
	b.subscribers[topic][sub] = struct{}{}

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		sub.Listen(b.ctx, b)
	}()

	log.Printf("Subscriber added to topic: %s\n", topic)

	return sub, nil
}

// Unsubscribe removes a subscriber
func (b *Broker) Unsubscribe(sub *Subscriber) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if subs, ok := b.subscribers[sub.Topic]; ok {
		if _, exists := subs[sub]; exists {
			delete(subs, sub)
			// Signal subscriber to stop
			close(sub.Messages)
			log.Printf("Subscriber removed from topic: %s\n", sub.Topic)
			if len(subs) == 0 {
				delete(b.subscribers, sub.Topic)
			}
		}
	}
}

// Publish sends a message to all subscribers of a topic
func (b *Broker) Publish(topic string, payload interface{}) error {
	b.lock.RLock()
	defer b.lock.RUnlock()

	subs, ok := b.subscribers[topic]
	if !ok {
		return errors.New("no subscribers for topic")
	}

	msg := Message{
		Topic:   topic,
		Payload: payload,
	}

	for sub := range subs {
		select {
		// Non-blocking publishing to prevent slow or unresponsive subscribers from blocking the entire system
		case sub.Messages <- msg:
		default:
			log.Printf("Subscriber channel full, skipping message for topic: %s\n", topic)
		}
	}

	log.Printf("Published message to topic: %s\n", topic)
	return nil

}

// Close gracefully shuts down the broker and all subscribers
func (b *Broker) Close() {
	b.cancel()
	b.wg.Wait()
	log.Println("Broker shut down gracefully.")
}
