package register

import (
	"how-dev-backend-test1/controller"
	"how-dev-backend-test1/system/middleware"
)

var c, m []any

// init 初始化
func init() {
	c = []any{
		controller.HttpOptionController{},
		controller.HealthCheckController{},
		controller.UserController{},
	}

	m = []any{
		middleware.CORS{},
	}
}

// Module 获取注册模块
func Module() []any {
	return append(c, m...)
}
