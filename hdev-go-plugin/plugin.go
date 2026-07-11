// Package hdevplugin 定义 How-Dev 框架的插件契约与全局注册表
//
// 该模块是框架中最底层的"扩展契约"模块，无任何外部依赖。
//
// 设计目的：
//   - 让应用启动器（hdev-go）与基础设施插件（hdev-gorm 等）通过本模块解耦
//   - 启动器只依赖本模块的接口去调度插件，不感知任何具体插件
//   - 插件只依赖本模块去注册自己，不感知启动器的存在
//
// 使用示例（插件方，例如 hdev-gorm）：
//
//	var enableOnce sync.Once
//	func EnableGorm() {
//	    enableOnce.Do(func() {
//	        hdevplugin.Register(myManager, "database.mysql")
//	    })
//	}
//
// 使用示例（启动器方，例如 hdev-go）：
//
//	for _, e := range hdevplugin.List() {
//	    _ = e.Plugin.Load(cfg.Get(e.ConfKey))
//	}
package hdevplugin

import (
	"sync"
)

// Plugin 框架插件接口
// 用于让基础设施模块（hdev-gorm、hdev-redis 等）在应用启动/停止时被自动调度
type Plugin interface {
	// Load 加载配置（从配置中心或本地配置文件读取）
	Load(conf interface{}) error
	// Reload 配置变更时重新加载
	Reload(conf interface{}) error
	// Setup 在所有 Process 启动前进行初始化
	Setup() error
	// Reset 重置插件状态（用于重启）
	Reset() error
	// Stop 收到关闭信号时调用
	Stop() error
	// Close 应用关闭时调用，释放资源
	Close() error
}

// Entry 插件注册条目（导出，供启动器遍历使用）
type Entry struct {
	// Plugin 插件实例
	Plugin Plugin
	// ConfKey 用于从配置中读取该插件的子配置（如 "database.mysql"）
	ConfKey string
}

// registry 全局插件注册表
type registry struct {
	mu      sync.RWMutex
	entries []*Entry
}

var globalRegistry = &registry{}

// Register 注册一个全局插件
// 推荐在插件包导出的 EnableXxx() 函数中调用（配合 sync.Once 保证幂等），
// 业务方在 main() 中显式调用 EnableXxx() 即可触发注册，例如：
//
//	var enableOnce sync.Once
//	func EnableGorm() {
//	    enableOnce.Do(func() {
//	        hdevplugin.Register(manager, "database.mysql")
//	    })
//	}
//
// p 为空时静默忽略；confKey 允许为空（表示插件不需要外部配置）
func Register(p Plugin, confKey string) {
	if p == nil {
		return
	}
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.entries = append(globalRegistry.entries, &Entry{
		Plugin:  p,
		ConfKey: confKey,
	})
}

// List 返回当前已注册的所有插件条目（返回拷贝，遍历安全）
// 供启动器（如 hdev-go）在生命周期各阶段遍历调度
func List() []*Entry {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	out := make([]*Entry, len(globalRegistry.entries))
	copy(out, globalRegistry.entries)
	return out
}

// Reset 清空注册表（仅测试场景使用）
func Reset() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.entries = nil
}

// -----------------------------------------------------------------------------
// ConfigAccessor：全局配置访问器
//
// 用于让插件在运行时（Load 阶段之外）也能读取任意配置项，同时不依赖具体的启动器。
// 启动器（如 hdev-go）在 App 初始化后调用 SetConfigAccessor 注入实现；
// 插件（如 hdev-gorm）通过 Config(key, callbacks...) 读取。
// -----------------------------------------------------------------------------

// ConfigAccessor 配置访问器契约
//
// 由启动器实现并通过 SetConfigAccessor 注入；插件通过 Config(...) 读取。
type ConfigAccessor interface {
	// Get 按点号分隔的 key 读取配置项；callbacks 在配置变更时被调用
	// 若配置未初始化或 key 不存在，返回 nil
	Get(key string, callbacks ...func(interface{})) interface{}
}

var (
	accessorMu     sync.RWMutex
	globalAccessor ConfigAccessor
)

// SetConfigAccessor 由启动器注入配置访问器实现
// 通常在 hdev-go.NewApp 中调用；多次调用以最后一次为准
func SetConfigAccessor(a ConfigAccessor) {
	accessorMu.Lock()
	defer accessorMu.Unlock()
	globalAccessor = a
}

// Config 读取全局配置（供插件使用）
// key 为点号分隔的配置路径；callbacks 在配置变更时被调用
// 当配置访问器未注入时返回 nil
func Config(key string, callbacks ...func(interface{})) interface{} {
	accessorMu.RLock()
	a := globalAccessor
	accessorMu.RUnlock()
	if a == nil {
		return nil
	}
	return a.Get(key, callbacks...)
}
