# hdev-go

How-Dev 框架的**应用启动器模块**，提供 Spring Boot 风格的链式启动 API：`NewApp().Process().Register().Run()`。

`hdev-go` 是**业务应用的直接依赖入口**——它把 `hdev-config`（配置加载）、`hdev-ioc`（IoC 容器）、`hdev-core`（接口）、`hdev-go-plugin`（插件契约）串联起来。基础设施插件（如 `hdev-gorm`）通过 [`hdev-go-plugin`](../hdev-go-plugin/README.md) 提供的契约自动接入，`hdev-go` 与具体插件互不感知。

## 🎯 功能特性

1. **App 启动器** — `NewApp / Process / Register / Run` 链式 API
2. **Plugin 生命周期调度** — 遍历 `hdev-go-plugin` 注册表，在启停期回调所有插件
3. **Process 调度** — 多个长期运行任务（HttpServer、Worker 等）按 yaml 配置启动
4. **配置桥接** — 加载主 yaml + 处理 `import`；把 `hdev-config.Config` 包装为 `hdev-core.PropertySource` 与 `hdev-ioc` 的 `configProvider`
5. **注入配置访问器** — `NewApp` 时把自身注入到 `hdev-go-plugin` 作为 `ConfigAccessor`，供插件读取运行时配置
6. **信号管理** — SIGINT / SIGTERM 优雅退出 + pid 文件
7. **全局访问** — `hdevgo.CurrentApp()` / `hdevgo.Config(key, callbacks...)`

## 🚀 快速开始

### 应用入口

```go
package main

import (
    hdevgo   "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
    hdevgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"    // 具名 import + 显式 EnableXxx()，替代旧的 blank-import
    http     "github.com/zhifenghao123/how-dev-go-framework/hdev-http-server"
    "how-dev-backend-test1/register"
)

func main() {
    hdevgorm.EnableGorm() // 显式启用 GORM 插件（注册到 hdev-go-plugin）

    hdevgo.NewApp(
        hdevgo.WithConf("conf/application.yml"),
        hdevgo.WithDebug(false),
    ).
        Process(http.HttpServer{}).             // 注册 Process（原型作用域）
        Register(register.Module()...).         // 注册业务 Bean（Controller/Service/Dao/Middleware）
        Run()                                   // 信号监听 → Plugin Setup → Process 调度 → 退出
}
```

### 主配置 `conf/application.yml`

```yaml
process:
  - class: HttpServer          # 对应 App.Process(...) 中注册的 Bean 类型名
    execute: Execute           # 调用的方法（默认 Execute）
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

## 🧩 App API

| 方法 | 说明 |
|------|------|
| `NewApp(opts...)` | 加载配置 → 初始化 BeanFactory → 遍历 `hdev-go-plugin` 触发 `Plugin.Load` → 注入 `ConfigAccessor` |
| `WithConf(file)` | 主配置文件路径，默认 `conf/application.yml` |
| `WithDebug(true)` | 注入 `debug=true` 到 `ctx` 中，供 BaseDao 等读取 |
| `.Process(beans...)` | 注册 Process Bean（**原型作用域**）；yaml `process:` 中 `class` 对应类型名 |
| `.Register(beans...)` | 注册业务 Bean（单例）；支持单个、多个、`[]interface{}` 切片 |
| `.Run()` | 监听信号 → Plugin.Setup → 调度 Process → 等待 → Plugin.Close |
| `.Close()` | 主动关闭（不等 Process 结束） |
| `.BeanFactory()` | 暴露内部 `*hdevioc.DefaultBeanFactory` |
| `.Config()` | 暴露 `*ConfigPropertySource` |
| `.Context()` | 暴露 `context.Context`，用于监听取消信号 |

### 全局访问

```go
hdevgo.CurrentApp()                                // 拿到当前 App
hdevgo.Config("database.mysqlLog.level", func(v interface{}) {
    // 配置变更回调
})
```

> 💡 **给插件开发者**：请使用 `hdevplugin.Config(...)` 而非 `hdevgo.Config(...)`，避免依赖 `hdev-go`。

## 🔌 Plugin 机制（契约在 `hdev-go-plugin`）

插件契约（`Plugin` 接口 + `Register` 函数 + `ConfigAccessor`）已抽离到独立模块 [`hdev-go-plugin`](../hdev-go-plugin/README.md)。`hdev-go` 只负责在应用生命周期各阶段**遍历注册表并回调**：

```
NewApp
  ├── loadPlugins()   → 遍历 hdevplugin.List()：Reset() → Load(subConfig)
  └── SetConfigAccessor(app)  → 供插件通过 hdevplugin.Config 读取配置

Run
  ├── setupPlugins()  → Setup()
  ├── (信号)
  ├── stopPlugins()   → Stop()
  └── closePlugins()  → Close()
```

插件方参考 [`hdev-go-plugin`](../hdev-go-plugin/README.md) 的示例，仅需依赖 `hdev-go-plugin`，无需依赖 `hdev-go`。

## 🧱 Process 模型

```yaml
process:
  - class: HttpServer         # Bean 名称（类型名）
    execute: Execute          # 调用的方法名（默认 Execute）
    params:                   # 反射注入到该 Bean 的字段
      Ip: 0.0.0.0
      Port: 8080
```

每个 Process 都会以独立 goroutine 调用对应 Bean 的 `Execute()`（或自定义方法）。

## 🔄 启动流程

```
NewApp(opts...)
  ├── LoadConfigFile(file)                    // 加载主配置 + 处理 import
  ├── loadPlugins(cfg)                        // 遍历 hdevplugin.List() 触发 Reset/Load
  ├── beanFactory = NewDefaultBeanFactory()
  └── setGlobalApp(a) → hdevplugin.SetConfigAccessor(app)

.Process(...).Register(...)                   // 把 Bean 写入 BeanFactory

.Run()
  ├── logPid()
  ├── go listenStopSignal()
  ├── setupPlugins()                          // Plugin.Setup()
  ├── for p in extractProcesses(): go runProcess(p)
  ├── wg.Wait()
  └── closePlugins() / removePid()
```

## 📁 文件结构

```
hdev-go/
├── app.go             # App 启动器（NewApp / Process / Register / Run / Close）
├── plugin.go          # Plugin 生命周期调度器（loadPlugins / setupPlugins / ...）
├── process.go         # Process 模型 + 调度器
├── signal.go          # 信号 / pid 管理
├── config_source.go   # ConfigPropertySource + LoadConfigFile
└── export.go          # CurrentApp / Config / appConfigAccessor 适配器
```

## 🔗 依赖

- `hdev-core`
- `hdev-config`
- `hdev-ioc`
- `hdev-go-plugin`（插件契约）

**不依赖** `hdev-context`（纯容器组件）、`hdev-gorm`（基础设施插件）等具体扩展；插件通过 `hdev-go-plugin` 反向接入。

## ⚙️ Go 版本

- Go 1.23+
