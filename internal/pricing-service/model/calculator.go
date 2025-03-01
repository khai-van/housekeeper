package model

import "housekeeper/api/pricing"

type JobRequire struct {
	StartDate    uint64
	RequiredHour int32
}

type JobPrice struct {
	Price    pricing.CurrencyValue
	Currency pricing.Currency
}
