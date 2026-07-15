package hdev_gorm

import (
	"gorm.io/gorm"
	"sync"
)

type mConfigs struct {
	sync.Map
}

func (m *mConfigs) get(key string) (c *dbConfig, ok bool) {
	v, ok := m.Load(key)
	if !ok {
		return
	}

	c, ok = v.(*dbConfig)
	return
}

func (m *mConfigs) loop(f func(key string, c *dbConfig) bool) {
	m.Range(func(k, v interface{}) bool {
		sk, ok := k.(string)
		if !ok {
			//ignore and continue
			return true
		}

		conf, ok := v.(*dbConfig)
		if !ok {
			//ignore and continue
			return true
		}

		return f(sk, conf)
	})
}

type mConnections struct {
	sync.Map
}

func (m *mConnections) get(key string) (conn *gorm.DB, ok bool) {
	v, ok := m.Load(key)
	if !ok {
		return
	}

	conn, ok = v.(*gorm.DB)
	return
}

func (m *mConnections) loop(f func(key string, c *gorm.DB) bool) {
	m.Range(func(k, v interface{}) bool {
		sk, ok := k.(string)
		if !ok {
			//ignore and continue
			return true
		}

		rc, ok := v.(*gorm.DB)
		if !ok {
			//ignore and continue
			return true
		}

		return f(sk, rc)
	})
}

type mCallbacks struct {
	sync.Map
}

func (m *mCallbacks) upsert(key string, callback ...func(*gorm.DB)) {
	if v, ok := m.get(key); !ok {
		m.Store(key, callback)
	} else {
		v = append(v, callback...)
		m.Store(key, v)
	}
}

func (m *mCallbacks) get(key string) (c []func(client *gorm.DB), ok bool) {
	v, ok := m.Load(key)
	if !ok {
		return
	}

	c, ok = v.([]func(client *gorm.DB))
	return
}
