package calculator

import (
	"context"
	"fmt"
	"housekeeper/api/pricing"
	"housekeeper/internal/pricing-service/model"
	"housekeeper/pkg/utils"
	"time"
)

type PricingCalculator struct {
	cfg *Config
}

func NewPricingCalculator(cfg *Config) (*PricingCalculator, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is null")
	}

	// missing default price
	if cfg.DefaultPrice.Price.Value == 0 {
		return nil, fmt.Errorf("default price is required in config")
	}

	// checking config fix date
	for date := range cfg.FixedDatePrices {
		_, err := utils.ParseStandardDate(date)
		if err != nil {
			return nil, fmt.Errorf("invalid start date in config: %w", err)
		}
	}

	// checking config rule date
	for _, rule := range cfg.RuleBasedRules {
		_, err := utils.ParseStandardDate(rule.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date in config: %w", err)
		}

		_, err = utils.ParseStandardDate(rule.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date in config: %w", err)
		}
	}

	return &PricingCalculator{
		cfg: cfg,
	}, nil
}

func (c *PricingCalculator) CalculatePrice(ctx context.Context, input model.JobRequire) (*model.JobPrice, error) {
	// missing required hour
	if input.RequiredHour < 1 {
		return nil, fmt.Errorf("invalid require hour input")
	}

	unitPrice := c.cfg.DefaultPrice

	// fixed date price
	startJobDate := time.Unix(int64(input.StartDate), 0)
	if price, exist := c.cfg.FixedDatePrices[startJobDate.Format("2006-01-02")]; exist {
		unitPrice = price
		goto RESULT
	}

	// rule base price
	for _, rule := range c.cfg.RuleBasedRules {
		startDate, _ := utils.ParseStandardDate(rule.StartDate) // already check err
		endDate, _ := utils.ParseStandardDate(rule.EndDate)     // already check err

		// job between range start and end
		if startDate.Before(startJobDate) && startJobDate.Before(endDate) {
			unitPrice = rule.Price
			break
		}
	}

RESULT:

	return &model.JobPrice{
		Price: pricing.CurrencyValue{
			Value:   unitPrice.Price.Value * int64(input.RequiredHour),
			Decimal: unitPrice.Price.Decimal,
		},
		Currency: unitPrice.Currency,
	}, nil
}
