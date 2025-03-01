package controller

import (
	"context"
	"fmt"
	"housekeeper/api/send"
	sendservice "housekeeper/internal/send-service"
	"housekeeper/internal/send-service/config"
	"housekeeper/internal/send-service/mock"
	"housekeeper/pkg/rabbitmqx"
)

type SendServer struct {
	send.UnimplementedSendServiceServer
	rabbitmqClient *rabbitmqx.RabbitMQClient

	svc *sendservice.SendService
}

func NewSendServer(cfg *config.Config) (*SendServer, error) {
	rabbitmqClient, err := rabbitmqx.NewRabbitMQClient(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	sendSvc, err := sendservice.NewSendService(
		rabbitmqClient,
		mock.NewMockEmployeeService(),
	)
	if err != nil {
		return nil, fmt.Errorf("create send service: %w", err)
	}

	return &SendServer{
		rabbitmqClient: rabbitmqClient,
		svc:            sendSvc,
	}, nil
}

func (s *SendServer) SendJob(ctx context.Context, req *send.SendJobRequest) (*send.SendJobResponse, error) {
	return s.svc.SendJob(ctx, req)
}

// Close closes the RabbitMQ connection when the service is shut down
func (s *SendServer) Close() {
	if s.rabbitmqClient != nil {
		s.rabbitmqClient.Close()
	}
}
