package sendservice

import (
	"context"
	"housekeeper/internal/send-service/model"
)

type EmployeeService interface {
	GetAvailableEmployee(context.Context, model.GetAvailableEmployeeRequest) ([]model.EmployeeInfo, error)
	GetEmployeeInfo(context.Context, []string) ([]model.EmployeeInfo, error)
}
