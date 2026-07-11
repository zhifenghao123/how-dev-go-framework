# hdev-go-plugin

How-Dev 框架的**扩展契约模块**——只包含 Plugin 接口与全局注册表，**零外部依赖**。

## 🎯 定位

`hdev-go-plugin` 是框架中最底层的"契约模块"，用来打破启动器（`hdev-go`）与基础设施插件（`hdev-gorm` 等）之间的相互依赖，实现真正的解耦。

```
                       ┌─────────────┐
                       │ hdev-go-plugin │  ← 契约（Plugin 接口 + 注册表 + 配置访问器）
                       └──────┬──────┘
              ┌───────────────┼──────────────┐
              ▼               ▼              ▼
        ┌──────────┐   ┌──────────┐   ┌──────────┐
        │ hdev-go  │   │hdev-gorm │   │hdev-redis│
        │  启动器  │   │  插件    │   │  插件    │  未来任意扩展
        └──────────┘   └──────────┘   └──────────┘
              ▲               ▲
              └───────┬───────┘
                      │
                 业务应用
```

**核心思想**：应用启动器与插件互不感知，都只依赖本模块的接口契约。

## 🚀 使用方式

### 插件方（提供扩展的框架，如 `hdev-gorm`）

推荐**显式启用**风格：用一个具名的 `EnableXxx()` 函数（`sync.Once` 保护）封装 `Register` 调用，业务方在 `main()` 里显式调用：

```go
package hdev_gorm

import (
    "log"
    "sync"

    hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"
)

var enableOnce sync.Once

// EnableGorm 显式启用插件（幂等）。
// 业务方：
//   func main() { hdev_gorm.EnableGorm(); hdevgo.NewApp(...).Run() }
func EnableGorm() {
    enableOnce.Do(func() {
        hdevplugin.Register(manager, "database.mysql")
        log.Printf("[hdev-plugin] enabled: hdev-gorm (database.mysql)")
    })
}
```

> ⚠️ **不推荐**直接在 `init()` 中 `Register` + 业务方 blank-import（`_ "..."`）：blank-import 容易被 IDE / goimports 误删，无启动日志，排查困难。

`Manager` 需要实现 `hdevplugin.Plugin` 接口：

```go
type Plugin interface {
    Load(conf interface{}) error       // 加载配置
    Reload(conf interface{}) error     // 配置热更新
    Setup() error                      // Process 启动前初始化
    Reset() error                      // 重置
    Stop() error                       // 停止信号回调
    Close() error                      // 关闭资源
}
```

### 启动器方（如 `hdev-go`）

```go
package hdevgo

import hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"

// 在启动流程各阶段遍历注册表回调
func setupPlugins() {
    for _, e := range hdevplugin.List() {
        _ = e.Plugin.Setup()
    }
}

// 在 App 初始化时注入配置访问器，供插件读取运行时配置
func setGlobalApp(a *App) {
    hdevplugin.SetConfigAccessor(&appConfigAccessor{app: a})
}
```

### 插件在运行时读取全局配置

```go
// hdev-gorm/connection.go
logLevel := hdevplugin.Config("database.mysqlLog.level", func(v interface{}) {
    // 配置变更回调
})
```

## 📋 API 一览

| 符号 | 说明 |
|------|------|
| `Plugin` interface | 6 方法生命周期契约（Load/Reload/Setup/Reset/Stop/Close） |
| `Entry` struct | 注册条目（含 Plugin 与 ConfKey） |
| `Register(p, confKey)` | 插件在 `EnableXxx()` 中调用（推荐 `sync.Once` 保护），登记自己 |
| `List()` | 启动器遍历所有已注册插件 |
| `Reset()` | 清空注册表（仅测试用） |
| `ConfigAccessor` interface | 配置访问器契约（Get 方法） |
| `SetConfigAccessor(a)` | 启动器在初始化时注入访问器实现 |
| `Config(key, callbacks...)` | 插件读取全局配置的入口 |

## 📁 文件结构

```
hdev-go-plugin/
├── plugin.go     # 全部实现（Plugin / Entry / Register / List / ConfigAccessor / Config）
└── go.mod        # 零依赖
```

## 🔗 依赖

**无**——本模块作为框架最底层的接口契约，不依赖任何其它 hdev-* 模块或第三方库。

## ⚙️ Go 版本

- Go 1.23+

## 💡 设计参考

- Go 惯用注册模式：`database/sql` + 各驱动 `init()` 注册
- 依赖倒置原则（DIP）：把接口抽到底层独立模块
- Spring Boot 的 `spring-boot-autoconfigure` 独立模块承担类似角色
