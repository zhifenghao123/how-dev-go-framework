# How-Dev Go Framework

一个基于 Go 语言的轻量级企业级应用框架，借鉴 Spring 框架的设计理念，提供 IoC 容器、配置管理、HTTP 服务、ORM 集成等核心能力。

## 📦 模块架构

```
how-dev-go-framework/
├── hdev-core/             # 核心接口与抽象（Environment / ApplicationContext / Event 等）
├── hdev-config/           # 配置管理（YAML / JSON / INI / Properties，多源、Profile、绑定）
├── hdev-ioc/              # IoC 容器（BeanDefinition / BeanFactory / 字段注入）
├── hdev-context/          # 应用上下文（DefaultApplicationContext / 事件多播 / Environment）
├── hdev-go-plugin/           # 插件契约（Plugin 接口 + 注册表 + ConfigAccessor，零依赖）
├── hdev-go/               # 应用启动器（App / Process / 信号）
├── hdev-http-server/      # HTTP 服务（路由 / 中间件 / 控制器）
├── hdev-gorm/             # GORM 数据访问层（作为 Plugin 通过 hdev-go-plugin 挂载）
├── how-dev-backend-test1/ # 完整示例工程 1
└── how-dev-backend-test2/ # 完整示例工程 2
```

### 模块依赖关系

```
                        ┌────────────┐
                        │  hdev-core │  核心接口（无外部依赖）
                        └─────┬──────┘
       ┌───────────────┬──────┼───────┬──────────────┐
       ▼               ▼      ▼       ▼              ▼
  ┌─────────┐  ┌──────────┐  │  ┌──────────┐  ┌──────────────┐
  │hdev-ioc │  │hdev-config│  │  │hdev-http │  │ hdev-context │  ← 纯上下文/容器
  │         │  │           │  │  │  -server │  │              │
  └────┬────┘  └────┬──────┘  │  └──────────┘  └──────────────┘
       │            │         │
       └──────┬─────┘         │
              ▼               │        ┌────────────┐
       ┌────────────┐         │        │hdev-go-plugin │  ← 插件契约（零依赖）
       │  hdev-go   │←────────┘        └──────┬─────┘
       │ (启动器)   │────────────────────────►│
       └────────────┘                         │
                                              │
                                     ┌────────┴─────────┐
                                     ▼                  ▼
                              ┌────────────┐    ┌──────────────┐
                              │ hdev-gorm  │    │ hdev-redis   │
                              │  (插件)    │    │ (未来扩展)    │
                              └────────────┘    └──────────────┘
```

> **核心解耦成果**：
> - `hdev-context` 是纯粹的应用上下文/容器（Bean 工厂、事件多播、Environment），仅依赖 `hdev-core`
> - `hdev-go-plugin` 是**扩展契约模块（零依赖）**——只有 Plugin 接口 + 全局注册表 + ConfigAccessor
> - `hdev-go` 是应用启动器，依赖 `hdev-config` + `hdev-ioc` + `hdev-core` + `hdev-go-plugin`；**业务应用直接依赖 `hdev-go`**
> - `hdev-gorm` 只依赖 `hdev-go-plugin`——**它不感知 `hdev-go` 的存在**
> - 未来任意扩展框架（`hdev-redis`、`hdev-mq` 等）只需依赖 `hdev-go-plugin` 即可无缝接入

### 模块职责一览

| 模块 | 职责 | 关键 API |
|------|------|----------|
| **hdev-core** | 接口与抽象层（无外部依赖） | `Environment`、`ApplicationContext`、`BeanFactory`、`ApplicationEvent` |
| **hdev-config** | 多格式配置加载、Profile、属性源、配置绑定 | `Init`、`ConfigurableEnvironment`、`ConfigurationBinder` |
| **hdev-ioc** | Bean 定义/工厂、字段反射注入（`wired`/`value`/`const` tag） | `DefaultBeanFactory`、`NewBeanDefinition` |
| **hdev-context** | 应用上下文与容器（Spring `ApplicationContext` 风格） | `DefaultApplicationContext`、`DefaultListableBeanFactory`、事件多播、`Environment` |
| **hdev-go-plugin** | 插件契约（零依赖）：Plugin 接口 + 注册表 + ConfigAccessor | `Plugin`、`Register`、`List`、`ConfigAccessor`、`Config` |
| **hdev-go** | **业务应用入口**：启动器、Process 调度、信号管理、生命周期调度插件 | `NewApp`、`WithConf`、`Process`、`Register`、`Run` |
| **hdev-http-server** | HTTP 服务、路由、中间件、控制器 | `HttpServer`、`IMiddleware`、`MiddleFunc` |
| **hdev-gorm** | GORM 封装，通过 `hdev-go-plugin` 自动挂载 | `BaseDao`、`Get(name)`、`SetDebug` |

## 🚀 快速开始

### 应用入口（test1 / test2 同款写法）

```go
package main

import (
    hdevgo   "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
    hdevgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
    http     "github.com/zhifenghao123/how-dev-go-framework/hdev-http-server"
    "how-dev-backend-test1/register"
)

func main() {
    hdevgorm.EnableGorm() // 显式启用 GORM 插件（注册到 hdev-go-plugin），幂等

    hdevgo.NewApp(hdevgo.WithConf("conf/application.yml")).
        Process(http.HttpServer{}).             // 注册 Process（按原型作用域）
        Register(register.Module()...).         // 注册业务 Bean（Controller/Service/Dao/Middleware）
        Run()                                   // 启动：信号监听 → Plugin Setup → Process 调度
}
```

### 主配置文件（`conf/application.yml`）

```yaml
process:
  - class: HttpServer
    params:
      Name: how-dev-go-app
      Ip: 0.0.0.0
      Port: 8080
      ReadTimeout: 1
      WriteTimeout: 30
      Route: conf/route.yml
import:
  - conf/database.yml @database
application:
  allowHost: localhost:3000,127.0.0.1:3000
```

### 字段注入用法（`wired` / `value` / `const`）

```go
type UserController struct {
    UserService *UserService `wired:"true"`           // 按类型注入
}

type CORS struct {
    Whitelist string `value:"application.allowHost"` // 从配置读取
}

type UserReadDao struct {
    *hdev_gorm.BaseDao `wired:"true"`                 // 嵌入 BaseDao 注入
}
func (d *UserReadDao) Database() string { return "hao_db3" }
```

## 🔧 启动流程

```
NewApp(opts...)                        // 位于 hdev-go
  │
  ├── 加载主配置文件（含 import 处理）
  ├── 初始化 BeanFactory（hdev-ioc.DefaultBeanFactory）
  ├── loadPlugins()             ← 遍历 hdevplugin.List() 触发已注册插件读取自身配置
  └── hdevplugin.SetConfigAccessor(app)  ← 供插件运行时读取全局配置
       │
       ▼
.Process(beans...).Register(beans...)
  │
  ▼
Run()
  ├── logPid() / 监听 SIGINT|SIGTERM
  ├── setupPlugins()            ← 遍历 hdevplugin.List() 调用 Plugin.Setup()
  ├── extractProcesses() → for { go runProcess(p) }
  ├── wait Processes done
  └── closePlugins() / removePid()
```

## 🔌 扩展新的插件

想开发一个新的基础设施插件（如 `hdev-redis`），只需：

1. **只依赖 `hdev-go-plugin`**（不依赖 `hdev-go`）
2. 实现 `hdevplugin.Plugin` 接口
3. 提供一个显式启用函数 `EnableRedis()`，内部使用 `sync.Once` 保护并调用 `hdevplugin.Register(myPlugin, "redis")`，同时打印一条启用日志（推荐格式：`[hdev-plugin] enabled: hdev-redis (redis)`）

业务方在 `main()` 中调用 `hdevredis.EnableRedis()` 即可自动生效——`hdev-go` 完全不需要改动。

> 📌 为什么不用 `init()` 匿名注册 + blank-import？
> - blank-import 容易被 IDE / goimports 误删，且没有启动日志、无法精确控制注册时机
> - 显式 `EnableXxx()` 语义清晰、可打印启用日志、可平滑演进为 Options 模式（如 `EnableXxx(WithConfKey("..."))`）

## 📚 模块文档

- [hdev-core](hdev-core/README.md) - 核心接口
- [hdev-config](hdev-config/README.md) - 配置管理
- [hdev-ioc](hdev-ioc/README.md) - IoC 容器
- [hdev-context](hdev-context/README.md) - 应用上下文（纯容器组件）
- [hdev-go-plugin](hdev-go-plugin/README.md) - **插件契约（新扩展模块的入口）**
- [hdev-go](hdev-go/README.md) - **应用启动器（业务入口）**
- [hdev-http-server](hdev-http-server/README.md) - HTTP 服务
- [hdev-gorm](hdev-gorm/README.md) - GORM 集成（Plugin 示例）
- [示例工程 test1 业务 API](how-dev-backend-test1/API_README.md)
- [示例工程 test2 业务 API](how-dev-backend-test2/API_README.md)

## 🛠️ 构建与测试

每个子模块都是独立的 Go module，使用 `replace` 指向本地路径：

```bash
# 构建框架本身
cd hdev-go-plugin && go build ./...
cd ../hdev-core && go build ./...
cd ../hdev-config && go build ./...
cd ../hdev-ioc && go build ./...
cd ../hdev-context && go build ./... && go test ./...
cd ../hdev-go && go build ./...
cd ../hdev-http-server && go build ./...
cd ../hdev-gorm && go build ./...

# 构建示例工程
cd ../how-dev-backend-test1 && go build -o /tmp/test1.bin ./main
cd ../how-dev-backend-test2 && go build -o /tmp/test2.bin ./main
```

## 🎯 设计要点

1. **轻量启动器** - 一行链式调用 `NewApp().Process().Register().Run()` 即可完成应用启动
2. **接口与实现分离** - `hdev-core` 仅定义接口，无外部依赖；其它模块按需实现
3. **上下文与启动解耦** - `hdev-context` 只负责容器语义，`hdev-go` 负责启动语义；两者职责独立
4. **启动器与插件解耦** - 通过 `hdev-go-plugin` 中间契约模块，`hdev-go` 与 `hdev-gorm` 互不感知
5. **Plugin 机制** - 基础设施模块通过显式 `EnableXxx()` 函数向 `hdev-go-plugin` 注入自己，业务工程在 `main()` 中显式启用；语义清晰、可幂等、有启动日志
6. **字段反射注入** - 通过 `wired` / `value` / `const` tag 完成 IoC，零侵入业务代码
7. **多 module 解耦** - 各模块独立 `go.mod`，可单独发版本

## 📄 许可证

MIT License
