package hdevInterface

import "context"

// IBeanFactory 是一个Bean工厂接口，也是控制反转容器接口，用于注册和获取实例
type IBeanFactory interface {
	// Register 注册一个实例到容器中
	Register(interface{})
	// Instance 获取一个实例,如果没有，则实例化并返回
	Instance(interface{}) interface{}
	// WithContext 设置上下文
	WithContext(ctx context.Context)
}
