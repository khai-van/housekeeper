package pricingservice

import (
	"context"
	"housekeeper/internal/pricing-service/model"
)

type PricingCalculator interface {
	CalculatePrice(context.Context, model.JobRequire) (*model.JobPrice, error)
}
