package hdev_gorm

import (
	"errors"
	"fmt"
	"github.com/goinggo/mapstructure"
)


// 数据库配置
type dbConfig struct {
	Driver      string
	Ip          string
	Port        int
	Database    string
	Username    string
	Password    string
	Charset     string
	MaxConn     int
	MaxIdle     int
	ConnMaxLift int
	DialTimeout int
	IsSecondary bool
}

// 数据库配置组装
func loadConfig(c interface{}) (map[string]*dbConfig, error) {
	if c == nil {
		return nil, errors.New("empty mysql config")
	}

	if conf, ok := c.(map[string]*dbConfig); ok {
		return conf, nil
	}

	if conf, ok := c.(map[string]interface{}); ok {
		config := make(map[string]*dbConfig)
		for k, v := range conf {
			config[k] = &dbConfig{}
			if err := mapstructure.Decode(v, config[k]); err != nil {
				return nil, errors.New(fmt.Sprintf("mysql [%s] config decode error : %v", k, err))
			}
		}
		return config, nil
	}

	return nil, errors.New("error mysql config")

}
