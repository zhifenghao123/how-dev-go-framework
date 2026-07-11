package hdev_gorm

import (
	"context"
	"log"

	"gorm.io/gorm"
)

// BaseDao 数据库基类
type BaseDao struct {
	*gorm.DB
	secondary *gorm.DB
	Ctx       context.Context `wired:"true"`
}

// Init 数据库连接获取
func (bd *BaseDao) Init(dao IDao) {
	if dao == nil || dao == bd {
		return
	}
	if dao.Database() == "" {
		log.Println("database of dao is not configured")
		return
	}

	if bd.Ctx != nil {
		debug = bd.Ctx.Value("debug").(bool)
	}

	conn := Get(dao.Database(), func(c *gorm.DB) {
		bd.DB = c
	})
	bd.DB = conn

	sKey := getSecondaryConfigKey()
	if sKey != "" {
		sConn := Get(sKey, func(c *gorm.DB) {
			bd.secondary = c
		})
		bd.secondary = sConn
	}
}

// Database 获取数据库名
func (bd *BaseDao) Database() string {
	return ""
}
