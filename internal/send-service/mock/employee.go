package mock

import (
	"context"
	"housekeeper/internal/send-service/model"
)

type MockEmployeeService struct{}

func NewMockEmployeeService() *MockEmployeeService {
	return &MockEmployeeService{}
}

func (*MockEmployeeService) GetAvailableEmployee(context.Context, model.GetAvailableEmployeeRequest) ([]model.EmployeeInfo, error) {
	return []model.EmployeeInfo{
		{
			ID:          "employee1",
			DeviceToken: "tokendevice1",
			DeviceType:  "android",
		},
		{
			ID:          "employee2",
			DeviceToken: "tokendevice2",
			DeviceType:  "android",
		},
	}, nil
}

func (*MockEmployeeService) GetEmployeeInfo(_ context.Context, ids []string) ([]model.EmployeeInfo, error) {
	res := []model.EmployeeInfo{}
	for _, id := range ids {
		res = append(res, model.EmployeeInfo{
			ID:          id,
			DeviceToken: "tokendevice" + id,
			DeviceType:  "android",
		})
	}

	return res, nil
}
