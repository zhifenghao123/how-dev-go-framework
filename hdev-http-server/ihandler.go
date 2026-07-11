package ihttp

import "net/http"

// IHandler interface of handler
type IHandler interface {
	Handle([]*action) *http.ServeMux
	Done()
}

// NewHandler get new handler
// execute all option of handler for init
func NewHandler(handlerMode int, h *HandlerConstruct) IHandler {
	if h.Handler == nil {
		h.Handler = DefaultHandler(handlerMode)
	}
	return func(handler IHandler) IHandler {
		for _, o := range h.Options {
			o(handler)
		}
		return handler
	}(h.Handler)
}

// HandlerConstruct construct for handler
type HandlerConstruct struct {
	Handler IHandler
	Options []IHandlerOption
}

func newHandler(handlerMode int, opt ...IHandlerOption) *HandlerConstruct {
	return &HandlerConstruct{
		Handler: DefaultHandler(handlerMode),
		Options: opt,
	}
}
