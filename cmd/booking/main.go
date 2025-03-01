package main

import (
	"fmt"
	"housekeeper/internal/booking-service/config"
	"housekeeper/internal/booking-service/controller"
	"housekeeper/pkg/utils"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configFile := os.Getenv("BOOKING_CONFIG_FILE") // e.g., "config.yaml"
	if configFile == "" {
		configFile = "config.yaml" // Default config file
	}

	cfg, err := utils.LoadConfig[config.Config](configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// init http server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	server, err := controller.NewBookingServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	server.RegisterRouter(e.Group(""))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
