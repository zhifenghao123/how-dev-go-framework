# hdev-context

How-Dev 框架的**应用上下文模块**（Spring `ApplicationContext` 风格），提供 Bean 工厂、事件多播器、`Environment` 实现，以及全套 Bean 生命周期能力。

> ⚠️ **应用启动器（`NewApp`/`Plugin`/`Process`/信号管理）已从本模块移出**，请使用 [`hdev-go`](../hdev-go/README.md) 作为业务应用的入口。`hdev-context` 现在是**纯上下文/容器组件**，不涉及配置文件加载、Process 调度等启动语义。

## 🎯 模块定位

| 提供 | 说明 |
|------|------|
| `DefaultApplicationContext` | 实现 `hdev-core.ConfigurableApplicationContext`，管理生命周期、事件、环境 |
| `DefaultListableBeanFactory` | 实现 `hdev-core.ConfigurableListableBeanFactory`，支持单例/原型作用域、后置处理器 |
| `SimpleApplicationEventMulticaster` | 事件发布/订阅 |
| `StandardEnvironment` + PropertySource 家族 | 属性源合并、Profile |

## 🚀 使用示例

### 独立使用应用上下文

```go
package main

import (
    hdevcontext "github.com/zhifenghao123/how-dev-go-framework/hdev-context"
    hdevcore    "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

type UserService struct {
    Name string
}

func (s *UserService) AfterPropertiesSet() error {
    // 初始化钩子
    return nil
}

func main() {
    ctx := hdevcontext.NewDefaultApplicationContext()

    // 注册 Bean 定义
    def := hdevcontext.NewBeanDefinition("userService", reflect.TypeOf((*UserService)(nil)))
    def.SetInitMethodName("AfterPropertiesSet")
    ctx.GetBeanFactory().RegisterBeanDefinition("userService", def)

    // 或直接注册单例实例
    // ctx.GetBeanFactory().RegisterSingleton("cfg", &Config{...})

    if err := ctx.Refresh(); err != nil {
        panic(err)
    }
    defer ctx.Close()

    userService := ctx.GetBean("userService").(*UserService)
    _ = userService
}
```

### 事件监听

```go
type StartupListener struct{}

func (l *StartupListener) OnApplicationEvent(event hdevcore.ApplicationEvent) {
    if _, ok := event.(*hdevcontext.ContextRefreshedEvent); ok {
        fmt.Println("Application context refreshed")
    }
}

ctx.AddApplicationListener(&StartupListener{})
```

### 事件多播器（可独立使用）

```go
mc := hdevcontext.NewSimpleApplicationEventMulticaster()
mc.AddApplicationListener(myListener)
mc.PublishEvent(hdevcore.NewGenericApplicationEvent("payload"))
```

### Environment 独立使用

```go
env := hdevcontext.NewStandardEnvironment()
env.AddPropertySource(hdevcontext.NewMapPropertySource("test", map[string]string{
    "app.name":    "demo",
    "app.version": "1.0",
}))
env.SetActiveProfiles("dev")
fmt.Println(env.GetProperty("app.name"))
```

## 🧩 关键类型

### `DefaultApplicationContext`

`hdev-core.ConfigurableApplicationContext` 的默认实现。

```go
type ConfigurableApplicationContext interface {
    GetBean(name string) interface{}
    GetBeanByType(beanType interface{}) interface{}
    ContainsBean(name string) bool
    IsSingleton(name string) bool
    GetEnvironment() Environment
    PublishEvent(event ApplicationEvent)
    SetEnvironment(env Environment)
    GetBeanFactory() ConfigurableListableBeanFactory
    Refresh() error
    Close() error
    IsActive() bool
}
```

### `DefaultListableBeanFactory`

- `RegisterSingleton(name, obj)` — 注册运行时单例
- `RegisterBeanDefinition(name, bd)` — 注册 Bean 定义（延迟或按需实例化）
- `PreInstantiateSingletons()` — 在 `Refresh()` 中被调用，预实例化非 lazy 的单例
- `AddBeanPostProcessor(...)` — Bean 生命周期扩展点

### `BeanDefinition`

```go
def := hdevcontext.NewBeanDefinition("myBean", reflect.TypeOf((*MyBean)(nil)))
def.SetScope(hdevcontext.ScopeSingleton) // 或 ScopePrototype
def.SetLazyInit(true)
def.SetInitMethodName("AfterPropertiesSet")
def.SetDestroyMethodName("Destroy")
```

## 📁 文件结构

```
hdev-context/
├── application_context.go            # DefaultApplicationContext / Builder
├── bean_factory.go                   # DefaultListableBeanFactory / BeanDefinition / BeanScope
├── application_event_multicaster.go  # SimpleApplicationEventMulticaster / 相关适配器
├── environment.go                    # StandardEnvironment / MapPropertySource / EnvironmentVariablesPropertySource
└── application_context_test.go       # 单元测试
```

## 🔗 依赖

- `hdev-core`（唯一依赖）

不再依赖 `hdev-config`、`hdev-ioc`、`hdev-go`。任何"启动应用"的场景请引入 [`hdev-go`](../hdev-go/README.md)。

## ⚙️ Go 版本

- Go 1.23+
