package healthcheck

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
