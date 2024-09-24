package main

import (
	"context"
	"log"
)

type Message struct {
	Topic   string
	Payload interface{}
}

// Subscriber represents a subscriber to a topic.
type Subscriber struct {
	Topic    string
	Messages chan Message
	Handler  func(Message)
}

func NewSubscriber(topic string, buffer int) *Subscriber {
	return &Subscriber{
		Topic:    topic,
		Messages: make(chan Message, buffer),
		Handler:  defaultHandler,
	}
}

// Listen starts listening for messages and handles them
func (s *Subscriber) Listen(ctx context.Context, broker *Broker) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled, stopping subscriber for topic: %s\n", s.Topic)
			return

		case msg, ok := <-s.Messages:
			if !ok {
				log.Printf("Message channel closed, stopping subscriber for topic: %s\n", s.Topic)
				return
			}
			s.HandleMessage(msg)
		}
	}

}

func (s *Subscriber) HandleMessage(msg Message) {
	if s.Handler != nil {
		s.Handler(msg)
	} else {
		defaultHandler(msg)
	}
}

func defaultHandler(msg Message) {
	log.Printf("Received message on topic '%s': %v\n", msg.Topic, msg.Payload)
}
