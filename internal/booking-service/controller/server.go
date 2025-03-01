package controller

import (
	"context"
	"fmt"
	bookingservice "housekeeper/internal/booking-service"
	"housekeeper/internal/booking-service/config"
	"housekeeper/internal/booking-service/model"
	"housekeeper/internal/booking-service/pricingsvc"
	"housekeeper/internal/booking-service/repository"
	"housekeeper/internal/booking-service/sendsvc"
	"housekeeper/pkg/mongox"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type BookingServer struct {
	svc           *bookingservice.BookingService
	pricingClient *pricingsvc.PricingClient
	sendClient    *sendsvc.SendClient
}

func NewBookingServer(cfg *config.Config) (*BookingServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Connect to MongoDB
	err := mongox.ConnectMongoDB(ctx, cfg.MongoDBURI, cfg.MongoDBName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MongoDB: %v", err)
	}

	// 2. Create gRPC Clients
	pricingClient, err := pricingsvc.NewPricingClient(cfg.PricingServiceAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Pricing Service client: %v", err)
	}

	sendClient, err := sendsvc.NewSendClient(cfg.SendServiceAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Send Service client: %v", err)
	}

	return &BookingServer{
		svc: bookingservice.NewBookingService(
			repository.NewRepository(),
			pricingClient,
			sendClient,
		),
		pricingClient: pricingClient,
		sendClient:    sendClient,
	}, nil
}

func (s *BookingServer) RegisterRouter(g *echo.Group) error {
	g.POST("/booking", s.CreateJobHandler)

	return nil
}

func (s *BookingServer) CreateJobHandler(ctx echo.Context) error {
	// parse request
	var request model.JobRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	// call svc
	jobs, err := s.svc.CreateJob(reqCtx, request)
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusOK, jobs)
}
