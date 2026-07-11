package ihttp

import (
	"context"
	"log"
	"net/http"
)

// TestMiddleware2 middleware for test
type TestMiddleware2 struct{}

// Handle handle for test middleware
func (TestMiddleware2) Handle(next MiddleFunc) MiddleFunc {
	return func(ctx context.Context, request *http.Request, response *Response) error {
		log.Println("enter2")
		r := next(ctx, request, response)
		response.SetHeader("test", "aaaaaaaaa")
		log.Println("exit2")
		return r
	}
}
