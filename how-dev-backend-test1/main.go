package main

import (
	"how-dev-backend-test1/register"

	hdevgo "github.com/zhifenghao123/how-dev-go-framework/hdev-go"
	hdevgorm "github.com/zhifenghao123/how-dev-go-framework/hdev-gorm"
	http "github.com/zhifenghao123/how-dev-go-framework/hdev-http-server"
)

func main() {
	hdevgorm.EnableGorm() // 显式启用 GORM 插件（等价于旧的 blank-import 副作用）

	hdevgo.NewApp(hdevgo.WithConf("conf/application.yml")).
		Process(http.HttpServer{}).
		Register(register.Module()...).
		Run()
}
