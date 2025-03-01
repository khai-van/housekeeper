package pricingservice

import (
	"context"
	"fmt"
	"housekeeper/api/pricing"
	"housekeeper/internal/pricing-service/calculator"
	"housekeeper/internal/pricing-service/model"
	"log"
)

type PricingService struct {
	calculator PricingCalculator
}

func NewPricingService() (*PricingService, error) {
	calculator, err := calculator.NewPricingCalculator(calculator.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("create calculator: %w", err)
	}
	return &PricingService{calculator: calculator}, nil
}

func (s *PricingService) GetPrice(ctx context.Context, req *pricing.GetPriceRequest) (*pricing.GetPriceResponse, error) {
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
