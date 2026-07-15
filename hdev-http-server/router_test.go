package ihttp

import (
	"reflect"
	"testing"
)

var router IRouter

func init() {
	router = NewRouter(&RouterConstruct{
		Router:  nil,
		Options: []IRouterOption{WithRouterFile("config/test_router.yml")},
	})
}

// test action
func Test_Action(t *testing.T) {
	actions := router.Action()
	m := make(map[string]*action, len(actions))
	for _, i := range actions {
		m[i.Method+"_"+i.Uri] = i
	}

	test := struct {
		want   map[string]*action
		result map[string]*action
	}{
		map[string]*action{
			"get_/test/pre1/pre2/test1": {
				Method:     "get",
				Uri:        "/test/pre1/pre2/test1",
				Controller: "TestController",
				Action:     "Test1",
				Middleware: []string{"TestMiddleware1", "TestMiddleware2", "TestMiddleware3"},
			},
			"get_/test/pre1/test2": {
				Method:     "get",
				Uri:        "/test/pre1/test2",
				Controller: "TestController",
				Action:     "Test1",
				Middleware: []string{"TestMiddleware1", "TestMiddleware2", "TestMiddleware4"},
			},
			"post_/test/pre1/test3": {
				Method:     "post",
				Uri:        "/test/pre1/test3",
				Controller: "TestController",
				Action:     "Test1",
				Middleware: []string{"TestMiddleware1", "TestMiddleware2"},
			},
			"get_/test/test4": {
				Method:     "get",
				Uri:        "/test/test4",
				Controller: "TestController",
				Action:     "Test2",
				Middleware: []string{"TestMiddleware1"},
			},
			"get_/test/test2/pre3/pre4/test8": {
				Method:     "get",
				Uri:        "/test/test2/pre3/pre4/test8",
				Controller: "TestController",
				Action:     "Test3",
				Middleware: []string{"TestMiddleware1", "TestMiddleware5", "TestMiddleware6", "TestMiddleware7"},
			},
			"get_/test/test2/test9": {
				Method:     "get",
				Uri:        "/test/test2/test9",
				Controller: "TestController",
				Action:     "Test4",
				Middleware: []string{"TestMiddleware1", "TestMiddleware5"},
			},
			"post_/test/test2/test10": {
				Method:     "post",
				Uri:        "/test/test2/test10",
				Controller: "TestController",
				Action:     "Test2",
				Middleware: []string{"TestMiddleware1", "TestMiddleware5"},
			},
		},
		m,
	}
	if !reflect.DeepEqual(test.want, test.result) {
		t.Errorf("Test_Action failed, expect %v, get %v ", test.want, test.result)
	}
}
