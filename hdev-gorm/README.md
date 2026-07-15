# hdev-gorm

How-Dev 框架的**数据访问层模块**，基于 GORM 封装，通过 [`hdev-go-plugin`](../hdev-go-plugin/README.md) 提供的插件契约自动挂载到应用启动器（如 [`hdev-go`](../hdev-go/README.md)）。

## 🎯 功能特性

- 🗄️ **多数据库支持**：MySQL、ClickHouse
- 🔄 **读写分离**：通过 `isSecondary: true` 标识从库
- 📊 **连接池管理**：可配置 `maxConn` / `maxIdle` / `connMaxLift`
- 📝 **DAO 模式**：嵌入 `BaseDao` 即可获得 `*gorm.DB`
- 🔍 **可配置 SQL 日志**：`silent / info / warn / error`
- 🔌 **Plugin 集成**：显式调用 `EnableGorm()` 注册到 `hdev-go-plugin`（不依赖 `hdev-go`）

## 🚀 快速开始

### 1. 在主入口显式调用 `EnableGorm()`

> ⚠️ **不再使用 blank-import 方式**（`_ "..."`）——旧的匿名 import 副作用容易被 IDE / goimports 误删，也没有启动日志。
> 现在改为**显式调用** `hdevgorm.EnableGorm()`：语义清晰、有启动日志、忘调不会静默失效。

```go
import (
    hdevgo   "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
    hdevgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
)

func main() {
    hdevgorm.EnableGorm() // ← 显式启用；幂等；打印启用日志

    hdevgo.NewApp(hdevgo.WithConf("conf/application.yml")).
        Process(/*...*/).
        Register(/*...*/).
        Run()
}
```

启动时会打印：

```
[hdev-plugin] enabled: hdev-gorm (database.mysql)
```

### 2. 主配置文件 `conf/application.yml` 引入数据库配置

```yaml
import:
  - conf/database.yml @database     # alias 为 "database"，配置键将变为 database.mysql...
```

### 3. 数据库配置 `conf/database.yml`

```yaml
mysql:
  hao_db3:                  # 数据库标识（DAO 通过 Database() 引用）
    driver: mysql           # mysql 或 clickhouse
    ip: 127.0.0.1
    port: 3306
    database: hao_db3
    username: root
    password: root_123
    charset: utf8
    maxConn: 20             # 最大连接数（默认 20）
    maxIdle: 10             # 最大空闲连接数（默认 10）
    connMaxLift: 300        # 连接生命周期（秒）
    dialTimeout: 5000       # 连接超时（毫秒），可选

mysqlLog:
  level: warn               # silent / info / warn / error
```

### 4. 编写 DAO

```go
package user_dao

import (
    "context"
    hgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
    "gorm.io/gorm"
)

type UserReadDao struct {
    *hgorm.BaseDao `wired:"true"`     // 通过 wired 注入 BaseDao（含 *gorm.DB）
}

// Database 返回该 DAO 绑定的数据库标识
func (d *UserReadDao) Database() string {
    return "hao_db3"
}

func (d *UserReadDao) GetUserByUserId(ctx context.Context, userId int64) (TableUser, error) {
    var user TableUser
    err := d.WithContext(ctx).
        Model(&TableUser{}).
        Where("user_id = ?", userId).
        First(&user).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return TableUser{}, nil
    }
    return user, err
}
```

### 5. 注册 DAO 到 IoC 容器

```go
// register/register.go
c = []any{
    &user_dao.UserReadDao{},
    &user_dao.UserWriteDao{},
    // ...
}
```

## 🔧 高级用法

### 读写分离

```yaml
mysql:
  primary:
    # 主库配置...
  reader-secondary:
    isSecondary: true        # 标识为从库
    ip: 127.0.0.2
    # ...
```

`BaseDao` 在 `Init` 时会自动找到 `isSecondary` 的从库连接。

### 直接获取连接

```go
db := hdev_gorm.Get("hao_db3")       // 直接拿到 *gorm.DB
```

### 调试模式

```go
hdev_gorm.SetDebug()                  // 开启 SQL 详细日志
```

或在 App 启动器中通过 `WithDebug(true)` 注入到 `ctx` 中，BaseDao 会读取并启用调试。

## 📦 Plugin 生命周期

`Manager` 实现了 `hdevplugin.Plugin` 接口：

| 方法 | 行为 |
|------|------|
| `Load(conf)` | 启动期解析配置并打开连接 |
| `Reload(conf)` | 配置变更时关闭旧连接并打开新连接 |
| `Setup()` | 占位（连接在 DAO `Init` 时按需创建） |
| `Reset()` | 关闭所有连接，清空状态 |
| `Stop()` | 收到关闭信号时占位 |
| `Close()` | 应用关闭时关闭所有连接 |

## 📁 文件结构

```
hdev-gorm/
├── init.go            # 包级变量的构造（manager / gLog），不再自动注册插件
├── enable.go          # EnableGorm() 显式启用入口（sync.Once 保护，幂等）
├── manager.go         # Manager 插件实现（Load/Reload/Reset/Close）
├── connection.go      # 连接构建、driver 选择、日志级别
├── config.go          # 配置结构体
├── base_dao.go        # BaseDao（含 *gorm.DB / secondary）
├── idao.go            # IDao 接口
├── callbacks.go       # GORM 回调（重试、慢查询）
├── map_wrapper.go     # 并发安全 map 包装
└── loggers/logger.go  # SQL 日志器
```

## 🔗 依赖

- `github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin`（插件契约，唯一的框架依赖）
- `gorm.io/gorm`、`gorm.io/driver/mysql`、`gorm.io/driver/clickhouse`

## ⚙️ Go 版本

- Go 1.23+
