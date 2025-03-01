package calculator_test

import (
	"context"
	"housekeeper/api/pricing"
	"housekeeper/internal/pricing/calculator"
	"housekeeper/internal/pricing/model"
	"reflect"
	"testing"
	"time"
)

func TestNewPricingCalculator(t *testing.T) {
	type args struct {
		cfg *calculator.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil config",
			args:    args{cfg: nil},
			wantErr: true,
		},
		{
			name: "missing default price",
			args: args{cfg: &calculator.Config{
				DefaultPrice: model.JobPrice{
					Price:    pricing.CurrencyValue{Value: 0, Decimal: 2},
					Currency: pricing.Currency_USD,
				},
			}},
			wantErr: true,
		},
		{
			name: "invalid fixed date",
			args: args{cfg: &calculator.Config{
				DefaultPrice: model.JobPrice{
					Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
					Currency: pricing.Currency_USD,
				},
				FixedDatePrices: map[string]model.JobPrice{
					"invalid-date": {Price: pricing.CurrencyValue{Value: 1500, Decimal: 2}, Currency: pricing.Currency_USD},
				},
			}},
			wantErr: true,
		},
		{
			name: "invalid rule start date",
			args: args{cfg: &calculator.Config{
				DefaultPrice: model.JobPrice{
					Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
					Currency: pricing.Currency_USD,
				},
				RuleBasedRules: []calculator.RuleConfig{
					{StartDate: "invalid-date", EndDate: "2025-03-10", Price: model.JobPrice{Price: pricing.CurrencyValue{Value: 2000, Decimal: 2}, Currency: pricing.Currency_USD}},
				},
			}},
			wantErr: true,
		},
		{
			name: "invalid rule end date",
			args: args{cfg: &calculator.Config{
				DefaultPrice: model.JobPrice{
					Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
					Currency: pricing.Currency_USD,
				},
				RuleBasedRules: []calculator.RuleConfig{
					{StartDate: "2025-03-01", EndDate: "invalid-date", Price: model.JobPrice{Price: pricing.CurrencyValue{Value: 2000, Decimal: 2}, Currency: pricing.Currency_USD}},
				},
			}},
			wantErr: true,
		},
		{
			name: "valid config",
			args: args{cfg: &calculator.Config{
				DefaultPrice: model.JobPrice{
					Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
					Currency: pricing.Currency_USD,
				},
				FixedDatePrices: map[string]model.JobPrice{
					"2025-03-01": {Price: pricing.CurrencyValue{Value: 1500, Decimal: 2}, Currency: pricing.Currency_USD},
				},
				RuleBasedRules: []calculator.RuleConfig{
					{StartDate: "2025-03-05", EndDate: "2025-03-10", Price: model.JobPrice{Price: pricing.CurrencyValue{Value: 2000, Decimal: 2}, Currency: pricing.Currency_USD}},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.NewPricingCalculator(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPricingCalculator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Errorf("NewPricingCalculator() = %v", got)
			}
		})
	}
}

func TestPricingCalculator_CalculatePrice(t *testing.T) {
	type fields struct {
		cfg *calculator.Config
	}
	type args struct {
		ctx   context.Context
		input model.JobRequire
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.JobPrice
		wantErr bool
	}{
		{
			name: "invalid required hour",
			fields: fields{
				cfg: &calculator.Config{
					DefaultPrice: model.JobPrice{
						Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
						Currency: pricing.Currency_USD,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: model.JobRequire{
					RequiredHour: 0,
					StartDate:    uint64(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC).Unix()),
				},
			},
			wantErr: true,
		},
		{
			name: "fixed date price",
			fields: fields{
				cfg: &calculator.Config{
					DefaultPrice: model.JobPrice{
						Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
						Currency: pricing.Currency_USD,
					},
					FixedDatePrices: map[string]model.JobPrice{
						"2025-03-01": {Price: pricing.CurrencyValue{Value: 1500, Decimal: 2}, Currency: pricing.Currency_USD},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: model.JobRequire{
					RequiredHour: 2,
					StartDate:    uint64(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC).Unix()),
				},
			},
			want: &model.JobPrice{
				Price: pricing.CurrencyValue{
					Value:   1500 * 2,
					Decimal: 2,
				},
				Currency: pricing.Currency_USD,
			},
		},
		{
			name: "rule based price",
			fields: fields{
				cfg: &calculator.Config{
					DefaultPrice: model.JobPrice{
						Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
						Currency: pricing.Currency_USD,
					},
					RuleBasedRules: []calculator.RuleConfig{
						{
							StartDate: "2025-03-05",
							EndDate:   "2025-03-10",
							Price:     model.JobPrice{Price: pricing.CurrencyValue{Value: 2000, Decimal: 2}, Currency: pricing.Currency_USD},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: model.JobRequire{
					RequiredHour: 3,
					StartDate:    uint64(time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC).Unix()),
				},
			},
			want: &model.JobPrice{
				Price: pricing.CurrencyValue{
					Value:   2000 * 3,
					Decimal: 2,
				},
				Currency: pricing.Currency_USD,
			},
		},
		{
			name: "default price",
			fields: fields{
				cfg: &calculator.Config{
					DefaultPrice: model.JobPrice{
						Price:    pricing.CurrencyValue{Value: 1000, Decimal: 2},
						Currency: pricing.Currency_USD,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: model.JobRequire{
					RequiredHour: 4,
					StartDate:    uint64(time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC).Unix()),
				},
			},
			want: &model.JobPrice{
				Price: pricing.CurrencyValue{
					Value:   1000 * 4,
					Decimal: 2,
				},
				Currency: pricing.Currency_USD,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := calculator.NewPricingCalculator(tt.fields.cfg)
			got, err := c.CalculatePrice(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PricingCalculator.CalculatePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PricingCalculator.CalculatePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
