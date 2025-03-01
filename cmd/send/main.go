package main

import (
	"fmt"
	"housekeeper/api/send"
	"housekeeper/internal/send-service/controller"
	"housekeeper/pkg/utils"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	// config initial
	configFile := os.Getenv("SEND_CONFIG_FILE") // e.g., "config.yaml"
	if configFile == "" {
		configFile = "config.yaml" // Default config file
	}

	cfg, err := utils.LoadConfig[Config](configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// setup service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s, err := controller.NewSendServer(&cfg.Config)
	if err != nil {
		log.Fatalf("Failed to create send API service server: %v", err)
	}

	grpcServer := grpc.NewServer()
	send.RegisterSendServiceServer(grpcServer, s)

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("Closing Send Service resources...")
		s.Close()
		log.Println("Shutdown complete.")
		os.Exit(0)
	}()

	log.Printf("Send service listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
