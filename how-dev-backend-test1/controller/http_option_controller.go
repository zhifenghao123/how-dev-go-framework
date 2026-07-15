package controller

import (
	"context"
	"net/http"
)

type HttpOptionController struct {
}

func (h *HttpOptionController) HttpOption(ctx context.Context, req *http.Request) (rsp any, err error) {
	return nil, nil
}
