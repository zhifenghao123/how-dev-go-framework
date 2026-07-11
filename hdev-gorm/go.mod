module github.com/zhifenghao123/how-dev-go-framework/hdev-gorm

go 1.23.2

require (
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin v0.0.0
	gorm.io/driver/clickhouse v0.3.2
	gorm.io/driver/mysql v1.3.4
	gorm.io/gorm v1.23.6
)

require (
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
)

replace github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin => ../hdev-go-plugin
