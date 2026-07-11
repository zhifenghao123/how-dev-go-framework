module how-dev-backend-test1

go 1.23.2

require github.com/zhifenghao123/how-dev-go-framework/hdev-go v0.0.0

require github.com/zhifenghao123/how-dev-go-framework/hdev-http-server v0.0.0

require (
	github.com/zhifenghao123/how-dev-go-framework/hdev-gorm v0.0.0
	gorm.io/gorm v1.23.6
)

require github.com/golang/protobuf v1.5.2 // indirect

require (
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/zhifenghao123/how-dev-go-framework/hdev-config v0.0.0 // indirect
	github.com/zhifenghao123/how-dev-go-framework/hdev-core v0.0.0 // indirect
	github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin v0.0.0 // indirect
	github.com/zhifenghao123/how-dev-go-framework/hdev-ioc v0.0.0 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/clickhouse v0.3.2 // indirect
	gorm.io/driver/mysql v1.3.4 // indirect
)

replace github.com/zhifenghao123/how-dev-go-framework/hdev-go => ../hdev-go

replace github.com/zhifenghao123/how-dev-go-framework/hdev-config => ../hdev-config

replace github.com/zhifenghao123/how-dev-go-framework/hdev-core => ../hdev-core

replace github.com/zhifenghao123/how-dev-go-framework/hdev-ioc => ../hdev-ioc

replace github.com/zhifenghao123/how-dev-go-framework/hdev-http-server => ../hdev-http-server

replace github.com/zhifenghao123/how-dev-go-framework/hdev-gorm => ../hdev-gorm

replace github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin => ../hdev-go-plugin
