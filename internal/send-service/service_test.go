package sendservice_test

import (
	"context"
	"fmt"
	"housekeeper/api/send"
	"housekeeper/internal/send-service/config"
	"housekeeper/internal/send-service/controller"
	"housekeeper/internal/send-service/worker"
	"housekeeper/pkg/rabbitmqx"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// integration test for send-service
const (
	PORT          = 8082
	NUMBER_WORKER = 2
)

var (
	s             *controller.SendServer
	w             *sendWorker // New
	rabbitMQURL   = "amqp://guest:guest@localhost:5672/"
	rabbitMQQueue = "test_job_queue"
)

type sendWorker struct { // New
	rabbitmqClient *rabbitmqx.RabbitMQClient
	workers        []*worker.SendWorker
}

func setup() {
	// 1. Setup RabbitMQ Client
	rabbitClient, err := rabbitmqx.NewRabbitMQClient(rabbitMQURL, rabbitMQQueue)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}

	// 2. Setup Send Service Server
	cfg := &config.Config{
		RabbitMQURL:   rabbitMQURL,
		RabbitMQQueue: rabbitMQQueue,
	}
	s, err = controller.NewSendServer(cfg) // Create the  service
	if err != nil {
		log.Fatalf("Failed to create SendAPIServiceServer: %v", err)
	}

	// 3. Initialize InMemory gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	send.RegisterSendServiceServer(grpcServer, s)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// 4. Setup Send Worker Service
	w = &sendWorker{ // Create the Worker
		rabbitmqClient: rabbitClient,
	}
	for i := 0; i < NUMBER_WORKER; i++ {
		sendW, err := worker.NewSendWorker(cfg) // Create the Dispatcher
		if err != nil {
			log.Fatalf("Failed to create SendWorker: %v", err)
		}
		sendW.Start() // Start consuming messages from RabbitMQ
		w.workers = append(w.workers, sendW)
	}
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	defer func() {
		if w != nil {
			w.rabbitmqClient.Close()
			for _, sendWorker := range w.workers {
				sendWorker.Close()
			}
			log.Println("Close Rabbitmq connect")
		}
		log.Println("Close connect done")
		s.Close()
		log.Println("Close api done")
	}()

	os.Exit(code)
}

func TestSendJobIntegration(t *testing.T) {
	// 1. Create a gRPC Client
	ctx := context.Background()
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", PORT),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err) // test no err when init conn
	defer conn.Close()

	client := send.NewSendServiceClient(conn)

	// 2. Send a SendJob Request\
	req := &send.SendJobRequest{
		JobId:          "jobID1",
		JobDescription: "JobDescription",
		JobAddress:     "JobAddress",
		StartDate:      uint64(time.Now().Unix()),
		RequiredHour:   1,
	}
	_, err = client.SendJob(ctx, req)
	require.NoError(t, err)

	// Give the worker some time to process the message
	time.Sleep(5 * time.Second)
}
