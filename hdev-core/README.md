# hdev-core

How-Dev 框架的**核心接口与抽象层**，仅定义接口和类型，**不包含具体实现**，无任何外部依赖。

## 📌 模块定位

`hdev-core` 是整个框架的最底层模块，它的职责仅限于：

- 定义贯穿框架的统一接口（`ApplicationContext`、`BeanFactory`、`Environment` 等）
- 提供基础的事件抽象（`ApplicationEvent`、`ApplicationListener`）
- 提供 Bean 生命周期回调接口（`InitializingBean`、`DisposableBean`、`BeanPostProcessor`）

具体实现散落在 `hdev-context`、`hdev-config`、`hdev-ioc` 等模块中。

## 🧩 主要接口

### 1. ApplicationContext（应用上下文）

```go
type ApplicationContext interface {
    GetBean(name string) interface{}
    GetBeanByType(beanType interface{}) interface{}
    GetBeansByType(beanType interface{}) []interface{}
    ContainsBean(name string) bool
    IsSingleton(name string) bool
    GetType(name string) interface{}
    GetAliases(name string) []string
    GetEnvironment() Environment
    PublishEvent(event ApplicationEvent)
}

type ConfigurableApplicationContext interface {
    ApplicationContext
    SetEnvironment(env Environment)
    GetBeanFactory() ConfigurableListableBeanFactory
    Refresh() error
    Close() error
    IsActive() bool
}
```

### 2. BeanFactory 体系

按职责分层：

```
BeanFactory                          (基础查询)
   ↑
ListableBeanFactory                  (+ 列举 / 计数)
   ↑
ConfigurableListableBeanFactory      (+ 注册 / 后置处理器)
```

### 3. Environment（环境配置）

```go
type Environment interface {
    GetProperty(key string) string
    GetPropertyWithDefault(key, defaultValue string) string
    ContainsProperty(key string) bool
    GetActiveProfiles() []string
    AcceptsProfiles(profiles ...string) bool
    // ...
}

type PropertySource interface {
    GetName() string
    GetProperty(key string) interface{}
    ContainsProperty(key string) bool
}
```

### 4. Bean 生命周期

```go
type InitializingBean interface {
    AfterPropertiesSet() error
}

type DisposableBean interface {
    Destroy() error
}

type BeanPostProcessor interface {
    PostProcessBeforeInitialization(bean interface{}, beanName string) (interface{}, error)
    PostProcessAfterInitialization(bean interface{}, beanName string) (interface{}, error)
}
```

### 5. 事件系统

```go
type ApplicationEvent interface {
    GetSource() interface{}
    GetTimestamp() int64
}

type ApplicationListener interface {
    OnApplicationEvent(event ApplicationEvent)
}

type SmartApplicationListener interface {
    ApplicationListener
    SupportsEventType(eventType reflect.Type) bool
    SupportsSourceType(sourceType reflect.Type) bool
}
```

## 📁 文件结构

```
hdev-core/
├── application_context.go   # ApplicationContext / BeanFactory / BeanDefinition 接口族
├── application_event.go     # ApplicationEvent / Listener / 通用事件
├── bean_lifecycle.go        # InitializingBean / DisposableBean / BeanPostProcessor
├── environment.go           # Environment / PropertySource
└── resource.go              # Resource 加载抽象
```

## 🔗 依赖关系

- **依赖**：无（仅使用 Go 标准库）
- **被依赖**：`hdev-config`、`hdev-ioc`、`hdev-context`

## ⚙️ Go 版本

- Go 1.23+

## 设计理念

1. **接口隔离** - 每个接口职责单一，便于扩展
2. **依赖倒置** - 高层模块不依赖低层模块，都依赖抽象
3. **开闭原则** - 对扩展开放，对修改关闭
4. **约定优于配置** - 提供合理的默认值，减少配置负担

## 扩展点

- `BeanPostProcessor` - 自定义Bean初始化逻辑
- `ApplicationListener` - 监听应用事件
- `PropertySource` - 自定义属性源
- `BeanFactoryPostProcessor` - 修改Bean定义

## 版本信息

- **当前版本**: 1.0.0
- **Go版本要求**: 1.23+

## 贡献指南

欢迎提交Issue和Pull Request来改进这个模块。