package ihttp

// IHandlerOption type of option for handler
type IHandlerOption func(IHandler)

// WithIoc ioc option for handler
func WithIoc(ioc IBeanFactory) IHandlerOption {
	return func(handler IHandler) {
		if h, ok := handler.(*Handler); ok {
			h.ioc = ioc
		}
	}
}
