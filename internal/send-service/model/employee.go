package model

type EmployeeInfo struct {
	ID          string `json:"id"`
	DeviceToken string `json:"deviceToken"` // token device
	DeviceType  string `json:"deviceType"`  // android, ios
}

type GetAvailableEmployeeRequest struct {
	JobAddress   string
	StartDate    uint64
	RequiredHour uint32
}

type WorkerMessage struct {
	EmployeeInfo `json:",inline"`

	JobId          string `json:"jobId"`
	JobDescription string `json:"jobDescription"`
	JobAddress     string `json:"jobAddress"`
	StartDate      uint64 `json:"startDate"`
	RequiredHour   uint32 `json:"requiredHour"`
}
