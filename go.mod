module github.com/zhifenghao123/how-dev-go-framework

go 1.23.2

require (
	github.com/zhifenghao123/how-dev-go-framework/hdev-config v0.0.0
	github.com/zhifenghao123/how-dev-go-framework/hdev-context v0.0.0-00010101000000-000000000000
	github.com/zhifenghao123/how-dev-go-framework/hdev-core v0.0.0
	github.com/zhifenghao123/how-dev-go-framework/hdev-ioc v0.0.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// 本地模块替换
replace github.com/zhifenghao123/how-dev-go-framework/hdev-core => ./hdev-core

replace github.com/zhifenghao123/how-dev-go-framework/hdev-context => ./hdev-context

replace github.com/zhifenghao123/how-dev-go-framework/hdev-config => ./hdev-config

replace github.com/zhifenghao123/how-dev-go-framework/hdev-ioc => ./hdev-ioc

replace github.com/zhifenghao123/how-dev-go-framework/hdev-go => ./hdev-go

replace github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin => ./hdev-go-plugin

replace github.com/zhifenghao123/how-dev-go-framework/hdev-gorm => ./hdev-gorm

replace github.com/zhifenghao123/how-dev-go-framework/hdev-http-server => ./hdev-http-server
