package domain

const ServiceName = "analytics-service"

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

func NewHealth() Health {
	return Health{
		Service: ServiceName,
		Status:  "ok",
	}
}
