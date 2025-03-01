package main

import (
	"fmt"
	"housekeeper/api/pricing"
	pricingservice "housekeeper/internal/pricing-service"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
)

const DEFAULT_PORT = 8081
const ENV_PORT_KEY = "PRICING_PORT"

func main() {
	// get port
	port := DEFAULT_PORT
	envPort := os.Getenv(ENV_PORT_KEY)
	if envPort != "" {
		parsedPort, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatalf("Invalid PRICING_PORT environment variable: %v", err)
		}
		port = parsedPort
	}

	// listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// init service
	s, err := pricingservice.NewPricingServiceServer()
	if err != nil {
		log.Fatalf("Failed to create pricing service server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pricing.RegisterPricingServiceServer(grpcServer, s)

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("Closing Pricing Service resources...")
		log.Println("Shutdown complete.")
		os.Exit(0)
	}()

	// start server
	log.Printf("----Pricing service start on %v----", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
