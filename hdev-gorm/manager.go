package hdev_gorm

import (
	"log"
	"reflect"

	hdevplugin "github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin"
	"gorm.io/gorm"
)

var _ hdevplugin.Plugin = (*Manager)(nil)

// Manager 数据库连接管理
type Manager struct {
	config      *mConfigs
	connections *mConnections
	callback    *mCallbacks
}

func (m *Manager) removeConfig(c map[string]*dbConfig) {
	if len(c) == 0 {
		return
	}

	for key := range c {
		log.Println("removing config:", key)
		if conn, ok := m.connections.get(key); ok {
			closeConnection(conn)
			m.connections.Delete(key)
		}
		m.config.Delete(key)
	}
}

// Load 数据库配置加载
func (m *Manager) Load(c interface{}) error {
	config, err := loadConfig(c)
	if err != nil {
		return err
	}

	//关闭被替换的链接
	var needClose []*gorm.DB
	defer func() {
		for _, conn := range needClose {
			closeConnection(conn)
		}
	}()

	for name, conf := range config {
		m.config.Store(name, conf)
		if conn, ok := m.connections.get(name); ok {
			needClose = append(needClose, conn)
			m.connections.Delete(name)
		}

		if cb, ok := m.callback.get(name); ok {
			newConn := Get(name)
			for _, b := range cb {
				b(newConn)
			}
		}
	}
	return nil
}

// Reload 数据库配置重载
func (m *Manager) Reload(c interface{}) error {
	config, err := loadConfig(c)
	if err != nil {
		return err
	}

	change := make(map[string]*dbConfig)
	remove := make(map[string]*dbConfig)
	for k, v := range config {
		if cv, ok := m.config.get(k); !ok {
			change[k] = v
			continue
		} else if !reflect.DeepEqual(v, cv) {
			change[k] = v
		}
	}

	m.config.loop(func(k string, value *dbConfig) bool {
		if _, ok := config[k]; !ok {
			remove[k] = value
		}
		return true
	})

	m.removeConfig(remove)
	return m.Load(change)
}

// Setup 数据库执行
func (m *Manager) Setup() error {
	return nil
}

// Reset 数据库连接池重置
func (m *Manager) Reset() error {
	_ = m.Close()

	m.connections.Clear()
	m.config.Clear()
	return nil
}

// Stop 插件停止
func (m *Manager) Stop() error {
	return nil
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	m.config.loop(func(key string, value *dbConfig) bool {
		return true
	})

	m.connections.loop(func(key string, value *gorm.DB) bool {
		closeConnection(value)
		return true
	})
	return nil
}

func (m *Manager) reconnect(name string) {
	log.Println("connection changed, reconnecting for:", name)

	var needClose []*gorm.DB
	defer func() {
		for _, conn := range needClose {
			closeConnection(conn)
		}
	}()

	if conn, ok := m.connections.get(name); ok {
		needClose = append(needClose, conn)
		m.connections.Delete(name)
	}

	if cb, ok := m.callback.get(name); ok {
		newConn := getWithInstance(name)
		for _, b := range cb {
			b(newConn)
		}
	}
}

// 数据库连接创建
// func (m *Manager) create(db string, inst *model.Instance) *gorm.DB {
func (m *Manager) create(db string) *gorm.DB {

	if conf, ok := m.config.get(db); ok {
		conn, err := connection(connOptions{
			name:   db,
			config: conf,
		})
		if err != nil {
			log.Println("connection create error : ", err)
			return nil
		}
		m.connections.Store(db, conn)
		log.Println("connection create success : ", db)
		return conn
	}
	log.Println("can't find config of database : ", db)
	return nil
}

// GetManager 获取数据库链接管理对象
func GetManager() *Manager {
	return manager
}

// Get 根据名称获取数据库链接
func Get(name string, callback ...func(*gorm.DB)) *gorm.DB {
	mutex.Lock()
	defer mutex.Unlock()

	if callback != nil {
		manager.callback.upsert(name, callback...)
	}

	if conn, ok := manager.connections.get(name); ok && conn != nil {
		return conn
	}
	return manager.create(name)
}

func getSecondaryDb() (conn *gorm.DB) {
	manager.config.loop(func(key string, value *dbConfig) bool {
		if value.IsSecondary {
			conn, _ = manager.connections.get(key)
			return false
		}
		return true
	})
	return
}

func getSecondaryConfigKey() (sKey string) {
	manager.config.loop(func(key string, value *dbConfig) bool {
		if value.IsSecondary {
			sKey = key
			return false
		}
		return true
	})
	return
}

func getWithInstance(name string) *gorm.DB {
	mutex.Lock()
	defer mutex.Unlock()

	return manager.create(name)
}
