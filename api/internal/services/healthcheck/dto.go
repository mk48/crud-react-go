package healthcheck

type HealthCheckResponse struct {
	Environment string `json:"environment"`
	Version     string `json:"version"`
}
