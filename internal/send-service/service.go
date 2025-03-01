package sendservice

import (
	"context"
	"fmt"
	"housekeeper/api/send"
	"housekeeper/internal/send-service/config"
	"housekeeper/internal/send-service/mock"
	"housekeeper/internal/send-service/model"
	"housekeeper/pkg/rabbitmqx"
	"log"
	"time"
)

type SendServiceServer struct {
	send.UnimplementedSendServiceServer
	rabbitmqClient *rabbitmqx.RabbitMQClient

	employeeSvc EmployeeService
}

// NewSendServiceServer creates a new SendServiceServer
func NewSendServiceServer(cfg *config.Config) (*SendServiceServer, error) {
	rabbitmqClient, err := rabbitmqx.NewRabbitMQClient(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	return &SendServiceServer{
		rabbitmqClient: rabbitmqClient,
		employeeSvc:    mock.NewMockEmployeeService(), // TODO: mock employee service for integration, complete it later
	}, nil
}

// SendJob receives a job ID and sends it to the dispatch workers
func (s *SendServiceServer) SendJob(ctx context.Context, req *send.SendJobRequest) (*send.SendJobResponse, error) {
	var listEmployeeInfo []model.EmployeeInfo
	var err error
	if len(req.EmployeeId) > 0 { // get specific data employee from specific id
		listEmployeeInfo, err = s.employeeSvc.GetEmployeeInfo(ctx, req.EmployeeId)
		if err != nil {
			return nil, fmt.Errorf("get employee info: %w", err)
		}
	} else { // get all available employee
		listEmployeeInfo, err = s.employeeSvc.GetAvailableEmployee(ctx, model.GetAvailableEmployeeRequest{
			JobAddress:   req.JobAddress,
			StartDate:    req.StartDate,
			RequiredHour: req.RequiredHour,
		})
		if err != nil {
			return nil, fmt.Errorf("get all available employee info: %w", err)
		}
	}

	// Publish the all message to worker for push message
	go s.pushToWorker(listEmployeeInfo, req)

	return &send.SendJobResponse{}, nil
}

func (s *SendServiceServer) pushToWorker(listEmployeeInfo []model.EmployeeInfo, req *send.SendJobRequest) {
	pCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, employee := range listEmployeeInfo {
		err := s.rabbitmqClient.Publish(pCtx, model.WorkerMessage{
			EmployeeInfo:   employee,
			JobId:          req.JobId,
			JobDescription: req.JobDescription,
			JobAddress:     req.JobAddress,
			StartDate:      req.StartDate,
			RequiredHour:   req.RequiredHour,
		})
		if err != nil {
			log.Println("failed to publish job to RabbitMQ: %w", err)
		}
	}
}

// Close closes the RabbitMQ connection when the service is shut down
func (s *SendServiceServer) Close() {
	if s.rabbitmqClient != nil {
		s.rabbitmqClient.Close()
	}
}
