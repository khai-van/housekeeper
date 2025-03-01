package worker

import (
	"encoding/json"
	"fmt"
	"housekeeper/internal/send-service/config"
	"housekeeper/internal/send-service/model"
	"housekeeper/pkg/rabbitmqx"
	"log"
	"time"
)

// SendWorker represents a worker that listens for jobs on a channel and dispatches them
type SendWorker struct {
	rabbitmqClient *rabbitmqx.RabbitMQClient
}

// NewDispatchWorker creates a new DispatchWorker
func NewSendWorker(cfg *config.Config) (*SendWorker, error) {
	rabbitmqClient, err := rabbitmqx.NewRabbitMQClient(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	return &SendWorker{
		rabbitmqClient: rabbitmqClient,
	}, nil
}

// Start starts the DispatchWorker, listening for jobs on the channel
func (w *SendWorker) Start() {
	go func() {
		err := w.rabbitmqClient.Consume(func(message []byte) error {
			log.Printf("Consuming message from queue: %s", message)
			return w.handle(message) // Dispatch the job
		})
		if err != nil {
			log.Printf("Error on Consumer: %s", err)
		}
	}()
}

func (w *SendWorker) Close() {
	if w.rabbitmqClient != nil {
		w.rabbitmqClient.Close()
	}
}

func (w *SendWorker) handle(message []byte) error {
	var data model.WorkerMessage
	if err := json.Unmarshal(message, &data); err != nil {
		return fmt.Errorf("unmarshal message err: %w", err)
	}

	// TODO: send to specific device of employee
	// Simulate the dispatching process
	time.Sleep(2 * time.Second) // Simulate some work

	return nil
}
