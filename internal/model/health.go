package model

type HealthMetrics struct {
	Database struct {
		Status  string `json:"status"`
		Error   string `json:"error"`
		Latency string `json:"latency"`
	} `json:"database"`
}

type ServiceMetrics struct {
	Status  string
	Error   string
	Latency string
}
