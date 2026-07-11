package ihttp

import (
	"context"
	"log"
	"net/http"
)

// TestMiddleware1 middleware for test
type TestMiddleware1 struct{}

// Handle handle for test middleware
func (TestMiddleware1) Handle(next MiddleFunc) MiddleFunc {
	return func(ctx context.Context, request *http.Request, response *Response) error {
		log.Println("enter")
		r := next(ctx, request, response)
		log.Println("exit")
		return r
	}
}
