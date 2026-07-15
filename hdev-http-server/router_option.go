package ihttp

// IRouterOption type of option for router
type IRouterOption func(IRouter)

// WithRouterFile option of router
func WithRouterFile(file string) IRouterOption {
	return func(router IRouter) {
		if r, ok := router.(*Router); ok {
			r.Init(file)
		}
	}
}
