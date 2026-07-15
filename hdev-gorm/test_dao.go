package hdev_gorm

import "fmt"

// Business 测试结构体：MCN业务
type Business struct {
	Id   int
	Name string
}

// TableName 指定关联数据表名
func (Business) TableName() string {
	return "business"
}

// TestDao test dao
type TestDao struct {
	BaseDao
}

// Database 获取数据库名称
func (td *TestDao) Database() string {
	return "hdev_hao"
}

// TestInsert 测试写入
func (td *TestDao) TestInsert() {
	business := &Business{
		Name: "aaaa",
	}
	td.Create(business)
	fmt.Println(business.Id)
}

// TestFind 测试查询
func (td *TestDao) TestFind() {
	business := &Business{}
	td.Select("id", "name").Where("id=?", 1).Find(business)
	fmt.Println(business.Id, business.Name)
}
