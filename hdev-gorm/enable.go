// Package hdev_gorm 的显式启用入口。
//
// 设计目的：
//   - 让应用显式调用 EnableGorm()，代替原来隐晦的 blank-import
//     （即 `_ "github.com/.../hdev-gorm"`），避免因遗漏 import 导致插件静默失效
//   - 显式调用可产生启动日志，便于运维/排查
//   - 未调用 EnableGorm() 时，本插件不会向 hdev-go-plugin 注册，
//     应用启动器 (hdev-go) 也就不会调度它——错误早暴露、意图更清晰
//
// 使用示例：
//
//	import (
//	    hdevgo   "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
//	    hdevgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
//	)
//
//	func main() {
//	    hdevgorm.EnableGorm()
//	    hdevgo.NewApp(hdevgo.WithConf("conf/application.yml")).
//	        Process(/*...*/).
//	        Register(/*...*/).
//	        Run()
//	}
package hdev_gorm

import (
	"log"
	"sync"

	hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"
)

var enableOnce sync.Once

// EnableGorm 显式启用 hdev-gorm 插件（幂等）。
//
// 调用时机：应在 hdevgo.NewApp(...).Run() 之前调用。
//
// 行为：
//  1. 向 hdev-go-plugin 全局注册表注册 Manager（confKey = "database.mysql"）
//  2. 打印一条启用日志
//
// 幂等：多次调用只会注册/打印一次。
func EnableGorm() {
	enableOnce.Do(func() {
		hdevplugin.Register(manager, "database.mysql")
		log.Printf("[hdev-plugin] enabled: hdev-gorm (database.mysql)")
	})
}
