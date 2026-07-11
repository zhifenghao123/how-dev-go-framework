package hdevgo

import (
	"log"

	hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"
)

// 本文件是应用启动器针对已注册插件的"生命周期调度器"。
//
// 插件契约（Plugin 接口 + 全局注册表 + Register 函数）已抽离到独立模块 hdev-go-plugin，
// 使得基础设施插件（hdev-gorm 等）无需依赖本模块 hdev-go，达成真正的解耦：
//
//	业务应用 ──→ hdev-go ────┐
//	                          ├──→ hdev-go-plugin  (只有接口契约)
//	业务应用 ──→ hdev-gorm ──┘
//
// 插件方使用（以 hdev-gorm 为例，推荐显式 EnableXxx() 方式）：
//
//	var enableOnce sync.Once
//	func EnableGorm() {
//	    enableOnce.Do(func() {
//	        hdevplugin.Register(manager, "database.mysql")
//	    })
//	}
//
// 启动器方（本文件）：
//
//	遍历 hdevplugin.List()，在应用生命周期各阶段回调插件

// loadPlugins 让所有插件读取配置
// 由 NewApp 调用；每个插件先 Reset，然后按 ConfKey 从主配置中抠出子配置喂给 Load
// 同时注册配置变更回调，配置热更新时会调用 plugin.Reload
func loadPlugins(cfg *ConfigPropertySource) {
	if cfg == nil {
		return
	}
	for _, e := range hdevplugin.List() {
		if err := e.Plugin.Reset(); err != nil {
			log.Printf("[hdev-go] plugin reset error: %v", err)
			continue
		}
		conf := cfg.Get(e.ConfKey, func(i interface{}) {
			if i == nil {
				return
			}
			if err := e.Plugin.Reload(i); err != nil {
				log.Printf("[hdev-go] plugin reload error: %v", err)
			}
		})
		if conf == nil {
			continue
		}
		if err := e.Plugin.Load(conf); err != nil {
			log.Printf("[hdev-go] plugin load error: %v", err)
		}
	}
}

// setupPlugins 调用所有插件的 Setup（Process 启动前）
func setupPlugins() {
	for _, e := range hdevplugin.List() {
		if err := e.Plugin.Setup(); err != nil {
			log.Printf("[hdev-go] plugin setup error: %v", err)
		}
	}
}

// stopPlugins 调用所有插件的 Stop（收到停止信号时）
func stopPlugins() {
	for _, e := range hdevplugin.List() {
		if err := e.Plugin.Stop(); err != nil {
			log.Printf("[hdev-go] plugin stop error: %v", err)
		}
	}
}

// closePlugins 调用所有插件的 Close（应用退出时）
func closePlugins() bool {
	for _, e := range hdevplugin.List() {
		if err := e.Plugin.Close(); err != nil {
			log.Printf("[hdev-go] plugin close error: %v", err)
		}
	}
	return true
}
