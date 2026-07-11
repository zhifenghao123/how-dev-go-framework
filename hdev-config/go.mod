module github.com/zhifenghao123/how-dev-go-framework/hdev-config

go 1.23.2

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/json-iterator/go v1.1.12
	github.com/zhifenghao123/how-dev-go-framework/hdev-core v0.0.0
	gopkg.in/ini.v1 v1.67.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/zhifenghao123/how-dev-go-framework/hdev-core => ../hdev-core