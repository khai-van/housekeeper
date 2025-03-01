package model

import (
	"housekeeper/api/pricing"
	"time"

	"github.com/kamva/mgm/v3"
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
