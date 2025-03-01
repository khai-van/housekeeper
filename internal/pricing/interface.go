package pricing

import (
	"context"
	"housekeeper/internal/pricing/model"
)

type PricingCalculator interface {
	CalculatePrice(ctx context.Context, input model.JobRequire) (*model.JobPrice, error)
}
