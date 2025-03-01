package controller

import (
	"context"
	"fmt"
	"housekeeper/api/pricing"
	pricingservice "housekeeper/internal/pricing-service"
)

type PricingServer struct {
	pricing.UnimplementedPricingServiceServer

	svc *pricingservice.PricingService
}

func NewPricingServer() (*PricingServer, error) {
	pricingSvc, err := pricingservice.NewPricingService()
	if err != nil {
		return nil, fmt.Errorf("create pricing service: %w", err)
	}

	return &PricingServer{
		svc: pricingSvc,
	}, nil
}

func (s *PricingServer) GetPrice(ctx context.Context, req *pricing.GetPriceRequest) (*pricing.GetPriceResponse, error) {
	return s.svc.GetPrice(ctx, req)
}
