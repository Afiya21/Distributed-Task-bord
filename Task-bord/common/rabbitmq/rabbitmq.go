package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RabbitClient handles the connection and channel for RabbitMQ
type RabbitClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// Event defines the standard structure for messages
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Connect initializes the RabbitMQ connection
func Connect(url string) (*RabbitClient, error) {
	var conn *amqp.Connection
	var err error

	// Retry logic for connection
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &RabbitClient{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close closes the connection and channel
func (rc *RabbitClient) Close() {
	if rc.Channel != nil {
		rc.Channel.Close()
	}
	if rc.Conn != nil {
		rc.Conn.Close()
	}
}

// Publish sends a message to a specific exchange
func (rc *RabbitClient) Publish(exchange, routingKey string, eventType string, payload interface{}) error {
	event := Event{
		Type:    eventType,
		Payload: payload,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return rc.Channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Consume starts consuming messages from a queue
func (rc *RabbitClient) Consume(queueName string, handler func(Event)) error {
	// Ensure queue exists
	q, err := rc.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	msgs, err := rc.Channel.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var event Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding event: %v", err)
				continue
			}
			handler(event)
		}
	}()

	return nil
}
