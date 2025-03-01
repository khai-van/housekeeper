package pricingservice

import (
	"context"
	"housekeeper/internal/pricing-service/model"
)

type PricingCalculator interface {
	CalculatePrice(ctx context.Context, input model.JobRequire) (*model.JobPrice, error)
}
