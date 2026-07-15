package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	hdevhttp "github.com/zhifenghao123/how-dev-go-framework/hdev-http-server"
)

// CORS 跨域中间件
type CORS struct {
	Whitelist string `value:"application.allowHost"`
}

var (
	MethodOptions = "options"

	AllowOrigin      = "Access-Control-Allow-Origin"
	AllowMethods     = "Access-Control-Allow-Methods"
	AllowCredentials = "Access-Control-Allow-Credentials"
	AllowHeaders     = "Access-Control-Allow-Headers"
	ExposeHeaders    = "Access-Control-Expose-Headers"

	AllowMethodsValue     = "POST, GET, OPTIONS, PUT, DELETE, PATCH"
	AllowCredentialsValue = "true"
	AllowHeadersValue     = "Content-Type, Accept, Origin, User-Agent, DNT, Cache-Control, X-Mx-ReqToken, X-Data-Type, X-Requested-With, Authorization"
	ExposeHeadersValue    = "Env"
)

// Handle 跨域处理
func (c CORS) Handle(next hdevhttp.MiddleFunc) hdevhttp.MiddleFunc {
	return func(ctx context.Context, request *http.Request, response *hdevhttp.Response) error {
		origin := request.Header.Get("Origin")

		println("cors middleware , HTTP_METHOD: %s, ORIGIN: %s, REFERER: %s", request.Method, origin, request.Referer())

		// 开发环境下允许所有源的跨域请求
		if origin != "" {
			response.SetHeader(AllowOrigin, origin)
		} else {
			// 如果没有Origin头，使用Referer作为备选
			referer := request.Referer()
			if referer != "" {
				path, err := url.Parse(referer)
				if err == nil {
					response.SetHeader(AllowOrigin, fmt.Sprintf("%s://%s", path.Scheme, path.Host))
				}
			}
		}

		response.SetHeader(AllowMethods, AllowMethodsValue)
		response.SetHeader(AllowCredentials, AllowCredentialsValue)
		response.SetHeader(AllowHeaders, AllowHeadersValue)
		response.SetHeader(ExposeHeaders, ExposeHeadersValue)

		println("cors middleware http method: %s", request.Method)
		if strings.ToLower(request.Method) == MethodOptions {
			return nil
		}
		return next(ctx, request, response)
	}
}
