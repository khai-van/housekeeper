package calculator

import (
	"housekeeper/api/pricing"
	"housekeeper/internal/pricing-service/model"
)

type Config struct {
	FixedDatePrices map[string]model.JobPrice // price in specific date
	RuleBasedRules  []RuleConfig              // price in range of rule date
	DefaultPrice    model.JobPrice            // default price
}

type RuleConfig struct {
	StartDate string
	EndDate   string
	Price     model.JobPrice
}

// TODO: extend to yaml | json config  or remote config
var config = Config{
	FixedDatePrices: map[string]model.JobPrice{},
	RuleBasedRules:  []RuleConfig{},
	DefaultPrice: model.JobPrice{
		Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
		Currency: pricing.Currency_USD,
	},
}

func GetConfig() *Config {
	return &config
}
