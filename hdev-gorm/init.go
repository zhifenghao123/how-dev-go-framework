package hdev_gorm

import (
	"log"
	"sync"
	"time"

	"github.com/zhifenghao123/how-dev-go-framework/hdev-gorm/loggers"
	"gorm.io/gorm/logger"
)

var (
	mutex   *sync.Mutex
	manager *Manager
	gLog    logger.Interface
	debug   bool
)

func init() {
	mutex = new(sync.Mutex)
	manager = &Manager{
		config:      &mConfigs{},
		connections: &mConnections{},
		callback:    &mCallbacks{},
	}
	gLog = loggers.New(log.Default(), loggers.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
	// 注意：不在此处向 hdev-go-plugin 注册插件；
	// 应用需显式调用 EnableGorm() 触发注册（见 enable.go）
}

// SetDebug 手动设置debug模式
func SetDebug() {
	debug = true
}
