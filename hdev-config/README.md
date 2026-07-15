# hdev-config

How-Dev 框架的**配置管理模块**，提供多格式配置文件加载、Profile 管理、属性源合并以及类型安全的配置绑定能力。

## 🎯 功能特性

- 多格式支持：**YAML / JSON / INI / Properties**
- 多属性源合并（文件、Map、环境变量）
- Profile 配置（`SPRING_PROFILES_ACTIVE`、`SetActiveProfiles`）
- 文件变更监听（基于 `fsnotify`）+ 回调式热更新
- 类型安全配置绑定（`ConfigurationBinder`）
- 实现 `hdev-core.Environment` 与 `PropertySource` 接口

## 🧩 主要类型

### 1. Config（核心入口）

实现 `hdev-interface.IConfig` 风格的统一访问 API：

```go
cfg, err := hdevconfig.Init(hdevconfig.Option{File: "conf/application.yml"})
// 链式取值
host := cfg.GetAsString("server.host")
port := cfg.GetAsInt("server.port")
list := cfg.GetAsArray("process")
m    := cfg.GetAsMap("database.mysql")

// 加载附加配置（带 alias 前缀）
_ = cfg.Load("conf/database.yml", "database")
```

支持的方法：`Get / GetAsString / GetAsInt / GetAsFloat / GetAsBool / GetAsArray / GetAsMap / GetAsStruct / GetAsStructArray`。

### 2. ConfigurableEnvironment（环境）

实现 `hdev-core.Environment`，按 PropertySource 优先级合并属性：

```go
env := hdevconfig.NewConfigurableEnvironment()
env.AddPropertySource(hdevconfig.NewMapPropertySource("app", map[string]interface{}{
    "server.port": "8080",
    "server.host": "localhost",
}))
fmt.Println(env.GetProperty("server.port")) // 8080
env.SetActiveProfiles("dev")
```

### 3. ConfigurationBinder（结构体绑定）

```go
type ServerConfig struct {
    Port int
    Host string
}

var sc ServerConfig
binder := hdevconfig.NewConfigurationBinder(env)
_ = binder.Bind("server", &sc)
```

### 4. PropertySource 实现

| 类型 | 用途 |
|------|------|
| `FilePropertySource` | YAML/JSON/INI/Properties 文件 |
| `MapPropertySource` | 内存 `map[string]interface{}` |
| `EnvironmentVariablesPropertySource` | OS 环境变量 |

### 5. ConfigFactory（快速创建）

```go
env, _ := hdevconfig.QuickStart()
env, _ = hdevconfig.QuickStartWithProfiles("dev", "local")
env, _ = hdevconfig.QuickStartWithFiles("config/app.yaml", "config/db.properties")
```

## 📁 文件结构

```
hdev-config/
├── config.go                   # Config（统一入口，watcher / Load / Get）
├── configurable_environment.go # ConfigurableEnvironment 实现
├── configuration_binder.go     # 类型安全配置绑定
├── config_factory.go           # ConfigFactory / QuickStart
├── file_property_source.go     # 文件属性源
├── value_extractor.go          # 类型转换提取器
├── ini_processor.go            # INI 处理器
├── json_processor.go           # JSON 处理器
├── yaml_processor.go           # YAML 处理器
├── loader.go                   # 加载流程
└── util.go                     # MapPropertySource / 工具
```

## 🔗 与其它模块协作

`hdev-context` 中 `App` 启动器加载主配置文件后，会把 `*hdevconfig.Config` 包装为 `ConfigPropertySource`：

- 既作为 `hdev-core.PropertySource` 暴露给 `Environment`
- 也作为 `configProvider` 注入给 `hdev-ioc.DefaultBeanFactory`，供 `value:"key"` 字段注入使用

## ⚙️ 依赖

- `github.com/zhifenghao123/how-dev-go-framework/hdev-core`
- `github.com/fsnotify/fsnotify`
- `github.com/json-iterator/go`
- `gopkg.in/yaml.v2`、`gopkg.in/ini.v1`、`github.com/goinggo/mapstructure`