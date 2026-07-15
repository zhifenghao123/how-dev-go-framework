package controller

import (
	"context"
	"how-dev-backend-test2/model"
	"net/http"
)

type HealthCheckController struct {
}

func (h *HealthCheckController) HealthCheck(ctx context.Context, req *http.Request) (rsp any, err error) {
	rsp = model.HealthCheckResult{
		HealthStatus: 1,
		Msg:          "account-security health check ok",
	}
	return
}
