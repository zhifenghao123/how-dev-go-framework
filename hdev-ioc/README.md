# hdev-ioc

How-Dev 框架的**依赖注入容器**，提供 BeanDefinition、BeanFactory、字段反射注入等核心能力。

## 🎯 功能特性

- **BeanDefinition** 元数据管理（作用域、依赖、生命周期方法）
- **DefaultBeanFactory** 默认 Bean 工厂实现
- **字段反射注入**：`wired` / `value` / `const` 三种 tag
- **作用域**：`singleton` / `prototype`
- **Bean 生命周期**：兼容 `Init(self)` 钩子（hdev-gorm BaseDao 风格）+ Spring 风格的 `AfterPropertiesSet` / `Destroy`
- **BeanPostProcessor**：初始化前后的扩展点
- 单例双重检查锁 + 提前曝光：支持简单循环依赖打破

## 🧩 核心 API

### 创建工厂 + 注册 Bean

```go
type UserService struct {
    Repo *UserRepository `wired:"true"`
}

type UserRepository struct {
    DBName string `value:"database.name"`
}

bf := hdevioc.NewDefaultBeanFactory()
bf.SetConfig(myConfig)            // 提供 Get(key) 的配置源
bf.WithContext(ctx)               // 关联 context.Context

bf.Register(&UserRepository{})    // 也可传值类型 UserRepository{}（自动转指针）
bf.Register(&UserService{})

us := bf.Instance(&UserService{}).(*UserService)
```

### 字段注入 tag 一览

| tag | 含义 | 示例 |
|-----|------|------|
| `wired:"true"` | 按字段类型自动注入 | `Repo *UserRepository \`wired:"true"\`` |
| `wired:"BeanName"` | 按 Bean 名称注入 | `Svc IService \`wired:"defaultSvc"\`` |
| `value:"key"` | 从配置源读取并赋值 | `Host string \`value:"server.host"\`` |
| `const:"text"` | 硬编码字符串常量 | `App string \`const:"my-app"\`` |

特殊预置类型：

- `context.Context`：通过 `wired:"true"` 自动注入工厂关联的 ctx
- 实现 `IBeanFactory` 接口的字段：自动注入工厂自身（用于 hdev-http-server 的 `HttpServer.Ioc` 字段）

### BeanDefinition 高级用法

```go
bd := hdevioc.NewBeanDefinition("MyBean", reflect.TypeOf((*MyBean)(nil)))
bd.SetScope(hdevioc.ScopePrototype) // 默认 singleton
bd.SetLazyInit(true)
bd.SetInitMethodName("AfterPropertiesSet")
bd.SetDestroyMethodName("Destroy")

bf.RegisterBeanDefinition("MyBean", bd)

// 预实例化所有非 lazy 的单例
_ = bf.PreInstantiateSingletons()
```

### BeanPostProcessor 扩展

```go
type LogProcessor struct{}

func (LogProcessor) PostProcessBeforeInitialization(bean interface{}, name string) (interface{}, error) {
    log.Printf("creating bean: %s", name)
    return bean, nil
}
func (LogProcessor) PostProcessAfterInitialization(bean interface{}, name string) (interface{}, error) {
    return bean, nil
}

bf.AddBeanPostProcessor(LogProcessor{})
```

## 📁 文件结构

```
hdev-ioc/
├── bean_definition.go     # BeanDefinition / Builder / Registry
├── bean_factory.go        # DefaultBeanFactory / 单例缓存 / 生命周期编排
├── field_injector.go      # 反射字段注入（wired/value/const）
└── annotation_scanner.go  # 实验性的注解扫描器（探索阶段）
```

## ♻️ Bean 生命周期

```
Register(bean) / RegisterBeanDefinition(name, bd)
       │
       ▼
GetBean(name) → getBean()
       │
       ├── singleton 缓存命中？─ 命中 ─▶ 直接返回
       │           ↓ 未命中
       │       创建空指针实例 → 加入单例缓存（提前曝光）
       │
       ├── PostProcessBeforeInitialization
       ├── injectFields()              ← wired / value / const
       ├── invokeInitHook()            ← Init(self) 兼容钩子
       ├── AfterPropertiesSet()        ← 实现 InitializingBean
       ├── PostProcessAfterInitialization
       └── 返回实例
```

## 🔗 依赖

- `github.com/zhifenghao123/how-dev-go-framework/hdev-core`

## ⚙️ Go 版本

- Go 1.23+
