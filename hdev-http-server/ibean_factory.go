package ihttp

// IBeanFactory hdev-http-server 内部使用的最小 IoC 容器接口
//
// 由调用方（如 hdev-context.App / hdev-ioc.DefaultBeanFactory）实现。
// 通过此接口解耦：hdev-http-server 不再依赖任何具体 IoC 实现，
// 任何符合 Instance(name string|reflect.Type|struct) interface{} 协议的容器均可注入。
type IBeanFactory interface {
	// Instance 通过名称、reflect.Type 或样例对象获取 Bean 实例
	// 实现方应支持以下三种入参形式：
	//   - string：按 Bean 名称查找
	//   - reflect.Type：按类型查找
	//   - struct/struct指针 实例：按类型查找
	Instance(bean interface{}) interface{}
}
