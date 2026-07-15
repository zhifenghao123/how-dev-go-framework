# hdev-http-server

How-Dev 框架的 **HTTP 服务模块**，提供路由配置、中间件链、控制器分发等 Web 能力。

## 🎯 功能特性

- 🌐 **HTTP Server**：基于 Go 标准库 `net/http`
- 🛣️ **路由管理**：YAML 配置驱动，支持嵌套分组
- 🔧 **中间件链**：可插拔，支持路由级 + 分组级 + Action 级
- 📋 **控制器分发**：通过反射调用 `Controller@Method`
- 🎛️ **两种处理模式**：
  - `0`（默认）：Path Mode —— 按 URI 直接路由
  - `1`：BodyAction Mode —— 单一 `/interface` 入口，根据 body 中 `Action` 字段分发
- 🔌 **零依赖**：仅本地定义 `IBeanFactory` 接口，与具体 IoC 容器解耦

## 🚀 快速开始

### 1. 在应用入口注册

```go
import (
    hdevgo "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
    http  "github.com/zhifenghao123/how-dev-go-framework/hdev-http-server"
)

func main() {
    hdevgo.NewApp(hdevgo.WithConf("conf/application.yml")).
        Process(http.HttpServer{}).             ← 把 HttpServer 作为 Process 注册
        Register(controllers...).               // 注册 Controller / Middleware（必须是 Bean）
        Run()
}
```

### 2. 主配置（`conf/application.yml`）

```yaml
process:
  - class: HttpServer
    params:
      Name: my-app
      Ip: 0.0.0.0
      Port: 8080
      ReadTimeout: 1            # 单位：秒
      WriteTimeout: 30
      Route: conf/route.yml     # 路由配置文件
      HandlerMode: 0            # 0=Path 模式 / 1=BodyAction 模式
```

### 3. 路由配置（`conf/route.yml`）

```yaml
groups:
  - prefix:
    middleware: CORS
    actions:
      POST /healthCheck:   { uses: HealthCheckController@HealthCheck }
      POST /addUser:       { uses: UserController@AddUser }
      POST /getUserById:   { uses: UserController@GetUserById }
      POST /login:         { uses: UserController@Login }
```

也支持嵌套分组与多个中间件（`|` 分隔）：

```yaml
groups:
  - prefix: /admin
    middleware: AuthMiddleware|LogMiddleware
    actions:
      GET /users: { uses: AdminController@ListUsers }
    groups:
      - prefix: /settings
        actions:
          GET /system: { uses: AdminController@SystemSettings }
```

## 🧩 Controller 编写

控制器只是普通的结构体 Bean，方法名对应路由配置中的 `@Method`：

```go
type UserController struct {
    UserService *service.UserService `wired:"true"`
}

func (u *UserController) Login(ctx context.Context, req *http.Request) (any, error) {
    var p model.LoginReq
    if err := tools.RequestParam(&p, req); err != nil {
        return nil, err
    }
    return u.UserService.Login(ctx, p)
}
```

### 支持的方法签名

控制器方法可使用以下任意签名：

```go
// 0 个参数
func (c *Ctrl) M()
// 1 个参数：context.Context
func (c *Ctrl) M(ctx context.Context)
// 2 个参数
func (c *Ctrl) M(ctx context.Context, req *http.Request)
// 3 个参数
func (c *Ctrl) M(ctx context.Context, req *http.Request, w http.ResponseWriter)
```

返回值支持：

```go
func (c *Ctrl) M(...)
func (c *Ctrl) M(...) error
func (c *Ctrl) M(...) interface{}
func (c *Ctrl) M(...) (interface{}, error)
```

## 🔧 Middleware 编写

实现 `IMiddleware` 接口：

```go
type CORS struct {
    Whitelist string `value:"application.allowHost"`
}

func (c CORS) Handle(next ihttp.MiddleFunc) ihttp.MiddleFunc {
    return func(ctx context.Context, req *http.Request, resp *ihttp.Response) error {
        // 前置处理
        resp.SetHeader("Access-Control-Allow-Origin", req.Header.Get("Origin"))
        if strings.ToLower(req.Method) == "options" {
            return nil
        }
        // 调用下一环
        return next(ctx, req, resp)
    }
}
```

`MiddleFunc` 定义：

```go
type MiddleFunc func(context.Context, *http.Request, *Response) error
```

## 📦 IBeanFactory 协议

`HttpServer.Ioc` 字段类型是模块**本地定义**的 `IBeanFactory`：

```go
type IBeanFactory interface {
    Instance(bean interface{}) interface{}
}
```

只要 IoC 容器实现了 `Instance(string|reflect.Type|struct) interface{}` 即可被注入。`hdev-ioc.DefaultBeanFactory` 已经满足。

## 📁 文件结构

```
hdev-http-server/
├── server.go              # HttpServer（Process 实现）
├── router.go              # 路由解析与构造
├── handler.go             # 分发器（Path / BodyAction 双模式）
├── ihandler.go            # IHandler 接口
├── irouter.go             # IRouter 接口
├── middware.go            # IMiddleware 接口
├── handler_option.go      # WithIoc 等 Option
├── router_option.go       # WithRouterFile 等 Option
├── ibean_factory.go       # 本地 IBeanFactory 接口
├── header.go              # 请求/响应头工具
├── json_content.go        # JSON 编解码
├── param.go               # 参数绑定（form / json / proto）
├── valid.go               # 参数校验（go-playground/validator）
├── response.go            # 响应结构
└── pb/                    # protobuf 测试桩
```

## 🔗 依赖

- `github.com/json-iterator/go`（JSON 编解码）
- `github.com/go-playground/validator/v10`（参数校验）
- `gopkg.in/yaml.v2`（路由配置解析）
- `google.golang.org/protobuf`（pb 响应支持）
- **不依赖** 框架其它模块（仅通过本地 `IBeanFactory` 接口与 IoC 容器解耦）

## ⚙️ Go 版本

- Go 1.23+
