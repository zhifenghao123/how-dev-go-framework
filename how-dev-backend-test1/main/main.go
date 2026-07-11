package main

import (
	hdevgo "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
	hdevgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
	http "github.com/zhifenghao123/how-dev-go-framework/hdev-http-server"
	"how-dev-backend-test1/register"
)

func main() {
	hdevgorm.EnableGorm() // 显式启用 GORM 插件

	hdevgo.NewApp(hdevgo.WithConf("conf/application.yml")).
		Process(http.HttpServer{}).
		Register(register.Module()...).
		Run()
}
