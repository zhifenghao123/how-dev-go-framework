package hdev_gorm

import (
	"fmt"

	hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"

	"log"
	"strings"
	"time"

	"github.com/zhifenghao123/how-dev-go-framework/hdev-gorm/loggers"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxOpenConn     = 20                  //最大连接数
	maxIdleConn     = 10                  //最大空闲数
	connMaxLifeTime = 14400 * time.Second //链接生命周期，4小时
	driverMysql     = "mysql"             //数据库驱动：mysql
	driverClick     = "clickhouse"        //数据库驱动：clickhouse
)

var (
	logLevel = logger.Warn
)

type connOptions struct {
	name   string
	config *dbConfig
}

// 获取数据库链接
func connection(co connOptions) (*gorm.DB, error) {
	c := co.config
	conn, err := source(co)
	if err != nil {
		return nil, err
	}

	db, err := conn.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(connMaxLifeTime)
	if c.MaxConn != 0 {
		db.SetMaxOpenConns(c.MaxConn)
	}
	if c.MaxIdle != 0 {
		db.SetMaxIdleConns(c.MaxIdle)
	}
	if c.ConnMaxLift != 0 {
		db.SetConnMaxLifetime(time.Duration(c.ConnMaxLift) * time.Second)
	}

	//备库不需要再切换
	if !co.config.IsSecondary {
		registerRetryCallbacks(conn)
	}
	return conn, nil
}

// 根据不同驱动创建链接
func source(co connOptions) (*gorm.DB, error) {
	logLevelConfig := hdevplugin.Config("database.mysqlLog.level", func(i interface{}) {
		SetLogLevel(i)
	})
	SetLogLevel(logLevelConfig)
	config := &gorm.Config{
		Logger: gLog,
	}
	if debug {
		config.Logger = gLog.LogMode(logger.Info)
	}

	c := co.config
	host, port := getMySqlHost(co)
	dialTimeout, hasTimeout := getTimeoutOpts(c)

	switch c.Driver {
	case driverMysql:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.Username, c.Password, host, port, c.Database, c.Charset)
		if hasTimeout {
			dsn += fmt.Sprintf("&timeout=%s", dialTimeout)
		}
		return gorm.Open(mysql.Open(dsn), config)
	case driverClick:
		dsn := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20",
			host, port, c.Username, c.Password, c.Database)
		if hasTimeout {
			dsn += fmt.Sprintf("&dial_timeout=%s", dialTimeout)
		}
		return gorm.Open(clickhouse.Open(dsn))
	default:
		return nil, nil
	}
}

func getTimeoutOpts(c *dbConfig) (dialTimeout time.Duration, ok bool) {
	//clickhouse默认的超时时间是30秒，mysql随系统设置，
	//这里为了向后兼容，默认还是不置顶，默认时间由驱动决定
	if c.DialTimeout > 0 {
		dialTimeout = time.Duration(c.DialTimeout) * time.Millisecond
		ok = true
	}
	return
}

func getMySqlHost(co connOptions) (host string, port int) {
	c := co.config

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[getHost]panic recovered, error:%v", r)
		}
		//兜底取配置参数
		if strings.EqualFold(host, "") {
			host = c.Ip
			port = c.Port
		}
	}()

	return
}

// 关闭数据库链接
func closeConnection(connection *gorm.DB) {
	if connection == nil {
		return
	}

	db, err := connection.DB()
	if err != nil {
		log.Println("mysql connection close error : ", err)
		return
	}
	defer func() {
		log.Println("closed mysql connection")
		connection = nil
	}()
	if err := db.Close(); err != nil {
		log.Println("mysql connection close error : ", err)
	}
}

func SetLogLevel(i interface{}) {
	logLevel = logger.Warn
	if config, ok := i.(string); ok {
		switch config {
		case "silent":
			logLevel = logger.Silent
		case "info":
			logLevel = logger.Info
		case "warn":
			logLevel = logger.Warn
		case "error":
			logLevel = logger.Error
		default:
			logLevel = logger.Warn
		}
	}
	gLog = loggers.New(log.Default(), loggers.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
}
