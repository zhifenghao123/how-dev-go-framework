package ihttp

import (
	"context"
	"log"
	"net/http"
)

type MiddleFunc func(context.Context, *http.Request, *Response) error

// IMiddleware middleware interface
type IMiddleware interface {
	Handle(MiddleFunc) MiddleFunc
}

// get list of middleware instance
func middleware(list []string, ioc IBeanFactory) (ml []IMiddleware) {
	l := len(list)
	if l == 0 {
		return
	}

	for i := l - 1; i >= 0; i-- {
		m := list[i]
		ins := ioc.Instance(m)
		if ins == nil {
			log.Printf("middleware [%s] instance failed", m)
			continue
		}

		if mi, ok := ins.(IMiddleware); ok {
			ml = append(ml, mi)
		} else {
			log.Printf("middleware [%s] is not impliment of IMiddleware", m)
		}
	}

	return
}