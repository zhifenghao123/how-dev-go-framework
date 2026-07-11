package ihttp

import (
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

const methodOptions = "options"
const (
	HandlerMode_Path       = 0
	HandlerMode_BodyAction = 1
)

var j jsoniter.API

func init() {
	j = jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}.Froze()
}

// reqActionBody 请求包
type reqActionBody struct {
	Action    string          `json:"Action"`
	Interface interfaceStruct `json:"interface"`
}

// interfaceStruct 接口关键信息
type interfaceStruct struct {
	Interface string      `json:"interface"`
	Param     interface{} `json:"param"`
}

// struct of dealers
type dealers struct {
	dealers map[string]map[string]*deal
}

// get dealers
func getDealers() *dealers {
	return &dealers{map[string]map[string]*deal{}}
}

// store dealer
func (d *dealers) store(uri, method string, dealer *deal) bool {
	if d, ok := d.dealers[uri]; ok {
		d[method] = dealer
		return true
	}
	d.dealers[uri] = map[string]*deal{method: dealer, methodOptions: dealer}
	return false
}

// find dealer
func (d *dealers) find(uri, method string) *deal {
	if dealer, ok := d.dealers[uri]; ok {
		if dd, ok := dealer[strings.ToLower(method)]; ok {
			return dd
		}
	}
	return nil
}

// struct of dealer
type deal struct {
	method     string
	execute    reflect.Value
	paramNum   int
	middleware []IMiddleware
}

type Handler struct {
	mode    int // 0:path, as well as default, 1:body action
	ioc     IBeanFactory
	dealers *dealers
	wg      sync.WaitGroup
}

// Handle register actions
func (h *Handler) Handle(actions []*action) *http.ServeMux {
	if h.ioc == nil {
		log.Fatalln("the ioc of default handler must be init")
	}
	h.dealers = getDealers()

	mux := http.NewServeMux()

	if h.mode == HandlerMode_Path {
		// path mode
		for _, a := range actions {
			h.wg.Add(1)
			go func(a *action) {
				defer h.wg.Done()

				dealer := h.dealer(a)
				if dealer == nil {
					return
				}

				if h.dealers.store(a.Uri, a.Method, dealer) {
					return
				}

				h.dealers.store(a.Uri, methodOptions, dealer)
				h.registerInPathMode(mux, a.Uri)
			}(a)
		}

		h.wg.Wait()
	} else {
		h.registerInBodyActionMode(mux, "/interface")
		for _, a := range actions {
			func(a *action) {
				dealer := h.dealer(a)
				if dealer == nil {
					return
				}
				if h.dealers.store(a.Uri, a.Method, dealer) {
					return
				}

			}(a)
		}
	}

	return mux
}

// Done check all request is deal
func (h *Handler) Done() {
	h.wg.Wait()
}

func (h *Handler) registerInPathMode(mux *http.ServeMux, uri string) {
	// Ensure URI starts with "/"
	if !strings.HasPrefix(uri, "/") {
		// 如果不是以"/"开头，则直接抛出异常
		log.Fatalln("invalid uri, must start with \"/\"")
	}
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		h.wg.Add(1)
		defer h.wg.Done()

		// init request params
		if err := request.ParseForm(); err != nil {
			log.Println("request parse error : ", err)
		}

		response := &Response{
			Header: make(map[string]string),
			Writer: writer,
		}

		defer output(writer, response)

		var (
			exec reflect.Value
			mid  []IMiddleware
			num  int
			ctx  context.Context
		)

		// find the dealer
		if dealer := h.dealers.find(uri, request.Method); dealer != nil {
			ctx, exec, mid, num = request.Context(), dealer.execute, dealer.middleware, dealer.paramNum
		}

		if !exec.IsValid() {
			response.Code, response.Data = http.StatusNotFound, "404 page not found"
			return
		}

		if err := execute(num, exec, mid, writer)(ctx, request, response); err != nil {
			response.Data = err.Error()
		}

	})
}
func (h *Handler) registerInBodyActionMode(mux *http.ServeMux, uri string) {
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		h.wg.Add(1)
		defer h.wg.Done()

		response := &Response{
			Header: make(map[string]string),
			Writer: writer,
		}
		defer output(writer, response)
		var (
			param  reqActionBody
			action string
		)
		if err := RequestParam(&param, true, false, request); err != nil {
			log.Printf("Cgw Route Error param=%v,error=%v \n", param, err.Error())
			return
		}
		action = param.Interface.Interface

		if action == "" {
			action = param.Action
		}

		if request.Method == "OPTIONS" && action == "" {
			action = "HttpOption"
		}

		// find the dealer
		if dealer := h.dealers.find(action, request.Method); dealer != nil {
			// init request params
			if err := request.ParseForm(); err != nil {
				log.Println("request parse error : ", err)
			}
			if err := execute(dealer.paramNum, dealer.execute, dealer.middleware, writer)(
				request.Context(), request, response); err != nil {
				response.Data = err.Error()
			}
			return
		}

		response.Code, response.Data = http.StatusNotFound, "404 action not found"
	})
}

// get dealer of action
func (h *Handler) dealer(a *action) *deal {
	con := reflect.ValueOf(h.ioc.Instance(a.Controller))
	if !con.IsValid() {
		log.Printf("please registed controller [%s]", a.Controller)
		return nil
	}

	function := con.MethodByName(a.Action)
	if !function.IsValid() {
		log.Printf("method [%s] in [%s] is invalid", a.Action, a.Controller)
		return nil
	}

	return &deal{
		method:     a.Method,
		execute:    function,
		paramNum:   function.Type().NumIn(),
		middleware: middleware(a.Middleware, h.ioc),
	}
}

func output(w http.ResponseWriter, r *Response) {
	for k, v := range r.Header {
		w.Header().Set(k, v)
	}

	var (
		byt []byte
		err error
	)
	defer func() {
		if err != nil {
			log.Println("io write string error : ", err)
		}
	}()

	if r.Code > 0 {
		w.WriteHeader(r.Code)
	}

	if str, ok := r.Data.(string); ok {
		_, err = io.WriteString(w, str)
		return
	}

	if msg, ok := r.Data.(proto.Message); ok {
		if byt, err = proto.Marshal(msg); err != nil {
			return
		}
		_, err = io.WriteString(w, string(byt))
		return
	}

	if byt, err = j.Marshal(r.Data); err != nil {
		return
	}
	_, err = io.WriteString(w, string(byt))
}

func params(ctx context.Context, n int, writer http.ResponseWriter, request *http.Request) (p []reflect.Value) {
	cx, rr, rw := reflect.ValueOf(ctx), reflect.ValueOf(request), reflect.ValueOf(writer)
	switch n {
	case 0:
		return
	case 1:
		return []reflect.Value{cx}
	case 2:
		return []reflect.Value{cx, rr}
	case 3:
		return []reflect.Value{cx, rr, rw}
	default:
		log.Fatalln("invalid num of action params, empty, " +
			"(*http.Request) or (*http.Request, http.ResponseWriter) is accepted")
	}
	return
}

func execute(num int, exec reflect.Value, mid []IMiddleware, writer http.ResponseWriter) MiddleFunc {
	f := func(ctx context.Context, request *http.Request, response *Response) error {
		res := exec.Call(params(ctx, num, writer, request))
		switch len(res) {
		case 0:
			return nil
		case 1:
			if err, ok := res[0].Interface().(error); ok {
				return err
			}
			response.Data = res[0].Interface()
			return nil
		case 2:
			response.Data = res[0].Interface()
			err := res[1].Interface()
			if err == nil {
				return nil
			}
			if e, ok := err.(error); ok {
				return e
			}
			return errors.New("invalid type of the second return element, error or nil accepted")
		default:
			return errors.New("invalid num of return, 2 or less is accepted")
		}
	}

	for _, i := range mid {
		f = i.Handle(f)
	}

	return f
}

// DefaultHandler get instance of default handler
func DefaultHandler(handlerMode int) IHandler {
	return &Handler{
		mode: handlerMode,
	}
}
