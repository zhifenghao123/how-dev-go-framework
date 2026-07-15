package ihttp

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

var i IBeanFactory

type TestIoc struct {
	Container map[string]reflect.Type
}

func (t *TestIoc) WithContext(ctx context.Context) {
	//TODO implement me
	//panic("implement me")
	println("with context")
}

func (t *TestIoc) Register(i interface{}) {
	if t.Container == nil {
		t.Container = make(map[string]reflect.Type)
	}
	r := reflect.TypeOf(i)
	t.Container[r.Name()] = r
}

func (t *TestIoc) Instance(i interface{}) interface{} {
	if s, ok := i.(string); ok {
		if r, ok := t.Container[s]; ok {
			return reflect.New(r).Interface()
		}
		return nil
	}

	return nil
}

func init() {
	ioc := new(TestIoc)
	ioc.Register(TestController{})
	ioc.Register(TestMiddleware1{})
	ioc.Register(TestMiddleware2{})
	i = ioc
}

// test of handler
func Test_Handler(t *testing.T) {
	h := NewHandler(0, &HandlerConstruct{
		Handler: nil,
		Options: []IHandlerOption{WithIoc(i)},
	})
	mux := h.Handle(router.Action())
	fmt.Println(mux)
}
