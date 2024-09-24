package main

import "log"

// Publisher represents a message publisher.
type Publisher struct {
	Broker *Broker
}

// NewPublisher creates a new Publisher.
func NewPublisher(broker *Broker) *Publisher {
	return &Publisher{
		Broker: broker,
	}
}

// Publish sends a message to a specified topic
func (p *Publisher) Publish(topic string, payload interface{}) {
	err := p.Broker.Publish(topic, payload)
	if err != nil {
		log.Printf("Failed to publish message to topic '%s': %v\n", topic, err)
	}
}
