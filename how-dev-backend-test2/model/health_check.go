package model

type HealthCheckResult struct {
	HealthStatus int    `json:"healthStatus"`
	Msg          string `json:"msg"`
}
