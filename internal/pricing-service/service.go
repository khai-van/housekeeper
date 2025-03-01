package pricingservice

import (
	"context"
	"fmt"
	"housekeeper/api/pricing"
	"housekeeper/internal/pricing-service/calculator"
	"housekeeper/internal/pricing-service/model"
	"log"
)

type PricingServiceServer struct {
	pricing.UnimplementedPricingServiceServer
	calculator PricingCalculator
}

func NewPricingServiceServer() (*PricingServiceServer, error) {
	calculator, err := calculator.NewPricingCalculator(calculator.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("create calculator: %w", err)
	}
	return &PricingServiceServer{calculator: calculator}, nil
}

func (s *PricingServiceServer) GetPrice(ctx context.Context, req *pricing.GetPriceRequest) (*pricing.GetPriceResponse, error) {
	price, err := s.calculator.CalculatePrice(ctx, model.JobRequire{
		StartDate:    req.StartDate,
		RequiredHour: int32(req.RequiredHour),
	})
	if err != nil {
		log.Printf("Error calculating price: %v", err)
		return nil, fmt.Errorf("calculating price: %v", err)
	}

	return &pricing.GetPriceResponse{
		Price:    &price.Price,
		Currency: price.Currency,
	}, nil
}
