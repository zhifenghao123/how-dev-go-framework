// Package hdevgo 提供应用启动器与全局配置访问入口
package hdevgo

import (
	"sync"

	hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"
)

// 全局 App 引用，便于业务代码通过 CurrentApp() 获取
var (
	globalAppMu sync.RWMutex
	globalApp   *App
)

// appConfigAccessor 把 App 适配为 hdevplugin.ConfigAccessor
// NewApp 会创建实例并通过 hdevplugin.SetConfigAccessor 注入，供插件读取全局配置
type appConfigAccessor struct {
	app *App
}

// Get 实现 hdevplugin.ConfigAccessor 接口
func (a *appConfigAccessor) Get(key string, callbacks ...func(interface{})) interface{} {
	if a.app == nil || a.app.cfg == nil {
		return nil
	}
	return a.app.cfg.Get(key, callbacks...)
}

// setGlobalApp 由 NewApp 内部调用：登记全局 App 引用 + 注入 ConfigAccessor 到 hdev-go-plugin
// 使得插件（hdev-gorm 等）可以通过 hdevplugin.Config(...) 读取运行时配置
func setGlobalApp(a *App) {
	globalAppMu.Lock()
	globalApp = a
	globalAppMu.Unlock()

	// 注入 ConfigAccessor 到 hdev-go-plugin，插件通过 hdevplugin.Config(...) 访问
	hdevplugin.SetConfigAccessor(&appConfigAccessor{app: a})
}

// CurrentApp 获取当前全局 App，若未初始化返回 nil
func CurrentApp() *App {
	globalAppMu.RLock()
	defer globalAppMu.RUnlock()
	return globalApp
}

// Config 全局配置访问入口（业务代码可用）
// key 为点号分隔的配置路径；callbacks 在配置变更时被调用
// 当全局 App 未初始化时返回 nil
//
// 注意：插件方（如 hdev-gorm）应改用 hdevplugin.Config，避免依赖 hdev-go
func Config(key string, callbacks ...func(interface{})) interface{} {
	globalAppMu.RLock()
	app := globalApp
	globalAppMu.RUnlock()
	if app == nil || app.cfg == nil {
		return nil
	}
	return app.cfg.Get(key, callbacks...)
}
