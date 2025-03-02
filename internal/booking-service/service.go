package bookingservice

import (
	"fmt"
	"housekeeper/api/send"
	"housekeeper/internal/booking-service/model"
	"log"
	"time"

	"golang.org/x/net/context"
)

type BookingService struct {
	repository BookingRepository

	pricingSvc PricingService
	sendSvc    SendService
}

func NewBookingService(
	repository BookingRepository,
	pricingSvc PricingService,
	sendSvc SendService,
) *BookingService {
	return &BookingService{
		repository: repository,
		pricingSvc: pricingSvc,
		sendSvc:    sendSvc,
	}
}

func (s *BookingService) CreateJob(ctx context.Context, req model.JobRequest) (*model.Job, error) {
	// validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 1. Call Pricing Service
	price, err := s.pricingSvc.GetPrice(ctx, uint64(req.StartDate), int32(req.RequiredHour))
	if err != nil {
		return nil, fmt.Errorf("get price err: %w", err)
	}

	// 2. Save Job to MongoDB
	job := model.Job{
		Description:  req.Description,
		CustomerID:   req.CustomerID,
		Address:      req.Address,
		StartDate:    time.Unix(req.StartDate, 0),
		RequiredHour: req.RequiredHour,
		Price:        *price.GetPrice(),
		Currency:     price.Currency,
	}
	if err = s.repository.CreateNewJob(ctx, &job); err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	// 3. Call Send Service
	go func() {
		mCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = s.sendSvc.SendJob(mCtx, &send.SendJobRequest{
			JobId:          job.ID.Hex(),
			JobDescription: job.Description,
			JobAddress:     job.Address,
			StartDate:      uint64(job.StartDate.Unix()),
			RequiredHour:   uint32(job.RequiredHour),
			EmployeeId:     req.EmployeeIDs,
		})
		if err != nil {
			log.Printf("Failed to send job: %v", err)
		}
	}()

	return &job, nil
}
