package hdev_gorm

import (
	"testing"
)

var (
	c       map[string]interface{}
	testDao *TestDao
)

func init() {
	c = map[string]interface{}{
		"hdev_hao": map[string]interface{}{
			"driver":      "mysql",
			"ip":          "",
			"port":        3306,
			"database":    "hdev_hao",
			"username":    "",
			"password":    "",
			"charset":     "utf8",
			"dialTimeout": 10,
		},
	}
	manager.Load(c)
	testDao = new(TestDao)
	testDao.Init(testDao)
}

// test find
func Test_Find(t *testing.T) {
	testDao.TestFind()
}

// test insert
//func Test_Insert(t *testing.T) {
//	testDao.TestInsert()
//}

// sql connection closed
func Test_Reload(t *testing.T) {
	testDao.TestFind()
	//c["hdev_hao"] = map[string]interface{}{
	//	"driver":   "mysql",
	//	"ip":       "",
	//	"port":     3306,
	//	"database": "hdev_hao",
	//	"username": "",
	//	"password": "",
	//	"charset":  "utf8",
	//}
	manager.Reload(c)
	testDao.TestFind()
}

// test
func Test_Reset(t *testing.T) {
	testDao.TestFind()
	manager.Reset()
	testDao.TestFind()
}

// test
func Test_Reset_Reload(t *testing.T) {
	testDao.TestFind()
	manager.Reset()
	manager.Reload(c)
	testDao.TestFind()
}
