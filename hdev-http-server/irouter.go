package ihttp

// struct of action
type action struct {
	Uri        string
	Method     string
	Controller string
	Action     string
	Middleware []string
}

// IRouter interface of router
type IRouter interface {
	Init(file string)
	Action() []*action
}

// NewRouter get new router
// execute all option of router for init
func NewRouter(r *RouterConstruct) IRouter {
	if r.Router == nil {
		r.Router = DefaultRouter()
	}
	return func(route IRouter) IRouter {
		for _, o := range r.Options {
			o(route)
		}
		return route
	}(r.Router)
}

// RouterConstruct construct of router
type RouterConstruct struct {
	Router  IRouter
	Options []IRouterOption
}

func newRouter(opt ...IRouterOption) *RouterConstruct {
	return &RouterConstruct{
		Router:  DefaultRouter(),
		Options: opt,
	}
}
