package bookingservice

import (
	"context"
	"housekeeper/api/pricing"
	"housekeeper/api/send"
	"housekeeper/internal/booking-service/model"
)

type (
	BookingRepository interface {
		CreateNewJob(ctx context.Context, job *model.Job) error
	}

	PricingService interface {
		GetPrice(ctx context.Context, startDate uint64, requiredHour int32) (*pricing.GetPriceResponse, error)
	}

	SendService interface {
		SendJob(ctx context.Context, req *send.SendJobRequest) error
	}
)
