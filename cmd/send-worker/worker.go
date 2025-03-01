package main

import (
	"housekeeper/internal/send-service/config"
	"housekeeper/internal/send-service/worker"
	"housekeeper/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFile := os.Getenv("SEND_WORKER_CONFIG_FILE") // e.g., "config.yaml"
	if configFile == "" {
		configFile = "config/config.yaml" // Default config file
	}

	cfg, err := utils.LoadConfig[config.Config](configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	s, err := worker.NewSendWorker(cfg)
	if err != nil {
		log.Fatalf("Failed to create send worker service: %v", err)
	}

	log.Println("Send worker service started")
	s.Start()

	// keep worker run
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan
	// Graceful shutdown
	log.Printf("Received signal: %v", sig)
	log.Println("Closing SendWorkerService resources...")
	s.Close() // Cast to access the Close method
	log.Println("Shutdown complete.")
}
