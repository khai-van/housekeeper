package model

import (
	"fmt"
	"housekeeper/api/pricing"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus int

const (
	JobPending JobStatus = iota
	JobReceived
	JobDone
)

type Job struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`

	Description string `json:"description" bson:"description"`
	CustomerID  string `json:"customerId" bson:"customerId"`
	Address     string `json:"address" bson:"address"`

	StartDate    time.Time `json:"startDate" bson:"startDate"`
	RequiredHour int       `json:"requiredHour" bson:"requiredHour"`

	Price    pricing.CurrencyValue `json:"price" bson:"price"`
	Currency pricing.Currency      `json:"currency" bson:"currency"`

	Status JobStatus `json:"status" bson:"status"`

	EmployeeID string `json:"employeeId,omitempty" bson:"employeeId,omitempty"`
}

type JobRequest struct {
	Description string `json:"description"`
	CustomerID  string `json:"customerId"`
	Address     string `json:"address"`

	StartDate    int64    `json:"startDate"`
	RequiredHour int      `json:"requiredHour"`
	EmployeeIDs  []string `json:"employeeIds"`
}

// Validate checks the JobRequest fields for valid data.
func (r *JobRequest) Validate() error {
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}

	if _, err := primitive.ObjectIDFromHex(r.CustomerID); err != nil {
		return fmt.Errorf("customerId is invalid: %w", err)
	}
	if r.Address == "" {
		return fmt.Errorf("address is required")
	}
	// Validate StartDate: must be a valid timestamp and not in the past (optional)
	if r.StartDate <= 0 {
		return fmt.Errorf("startDate must be a valid Unix timestamp")
	}
	if time.Unix(r.StartDate, 0).Before(time.Now()) {
		return fmt.Errorf("startDate cannot be in the past")
	}
	if r.RequiredHour < 1 {
		return fmt.Errorf("requiredHour must be at least 1")
	}

	return nil
}
