package rabbitmqx

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQClient(url string, queue string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare queue
	_, err = channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitMQClient{
		conn:      conn,
		channel:   channel,
		queueName: queue,
	}, nil
}

func (r *RabbitMQClient) Publish(ctx context.Context, data any) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal json data err: %w", err)
	}

	err = r.channel.PublishWithContext(ctx,
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

func (r *RabbitMQClient) Consume(handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := handler(d.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
		log.Println("Consumer stopped")
	}()

	return nil
}

func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
