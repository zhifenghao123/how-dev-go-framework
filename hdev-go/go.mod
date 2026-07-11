module github.com/zhifenghao123/how-dev-go-framework/hdev-go

go 1.23.2

require (
	github.com/zhifenghao123/how-dev-go-framework/hdev-config v0.0.0
	github.com/zhifenghao123/how-dev-go-framework/hdev-core v0.0.0
	github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin v0.0.0
	github.com/zhifenghao123/how-dev-go-framework/hdev-ioc v0.0.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/zhifenghao123/how-dev-go-framework/hdev-core => ../hdev-core

replace github.com/zhifenghao123/how-dev-go-framework/hdev-ioc => ../hdev-ioc

replace github.com/zhifenghao123/how-dev-go-framework/hdev-config => ../hdev-config

replace github.com/zhifenghao123/how-dev-go-framework/hdev-go-plugin => ../hdev-go-plugin
