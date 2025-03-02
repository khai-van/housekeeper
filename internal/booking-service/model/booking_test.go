package model_test

import (
	"housekeeper/internal/booking-service/model"
	"strings"
	"testing"
	"time"
)

func TestJobRequest_Validate(t *testing.T) {
	// Use a valid ObjectID hex string.
	validCustomerID := "507f1f77bcf86cd799439011"
	// Set up a future time and a past time.
	futureTime := time.Now().Add(time.Hour).Unix()
	pastTime := time.Now().Add(-time.Hour).Unix()

	tests := []struct {
		name    string
		req     model.JobRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "empty description",
			req: model.JobRequest{
				Description:  "",
				CustomerID:   validCustomerID,
				Address:      "123 Main St",
				StartDate:    futureTime,
				RequiredHour: 1,
			},
			wantErr: true,
			errMsg:  "description is required",
		},
		{
			name: "invalid customerId",
			req: model.JobRequest{
				Description:  "Clean the house",
				CustomerID:   "invalid", // not a valid ObjectID hex
				Address:      "123 Main St",
				StartDate:    futureTime,
				RequiredHour: 1,
			},
			wantErr: true,
			errMsg:  "customerId is invalid",
		},
		{
			name: "empty address",
			req: model.JobRequest{
				Description:  "Clean the house",
				CustomerID:   validCustomerID,
				Address:      "",
				StartDate:    futureTime,
				RequiredHour: 1,
			},
			wantErr: true,
			errMsg:  "address is required",
		},
		{
			name: "invalid startDate (zero)",
			req: model.JobRequest{
				Description:  "Clean the house",
				CustomerID:   validCustomerID,
				Address:      "123 Main St",
				StartDate:    0,
				RequiredHour: 1,
			},
			wantErr: true,
			errMsg:  "startDate must be a valid Unix timestamp",
		},
		{
			name: "startDate in the past",
			req: model.JobRequest{
				Description:  "Clean the house",
				CustomerID:   validCustomerID,
				Address:      "123 Main St",
				StartDate:    pastTime,
				RequiredHour: 1,
			},
			wantErr: true,
			errMsg:  "startDate cannot be in the past",
		},
		{
			name: "invalid requiredHour",
			req: model.JobRequest{
				Description:  "Clean the house",
				CustomerID:   validCustomerID,
				Address:      "123 Main St",
				StartDate:    futureTime,
				RequiredHour: 0,
			},
			wantErr: true,
			errMsg:  "requiredHour must be at least 1",
		},
		{
			name: "valid job request",
			req: model.JobRequest{
				Description:  "Clean the house",
				CustomerID:   validCustomerID,
				Address:      "123 Main St",
				StartDate:    futureTime,
				RequiredHour: 2,
			},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, want error containing %q", err.Error(), tt.errMsg)
			}
		})
	}
}
