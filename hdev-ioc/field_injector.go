package hdev_ioc

import (
	"context"
	"log"
	"reflect"

	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// 注入用的tag常量
const (
	TagWired = "wired"
	TagValue = "value"
	TagConst = "const"
)

// configProvider 配置提供者接口（弱依赖，避免循环引用 hdev-interface）
// 任何实现了 Get(key, callbacks...) 方法的对象都可作为配置源
type configProvider interface {
	Get(string, ...func(i interface{})) interface{}
}

// configurableFactory 表示工厂可关联配置和上下文
type configurableFactory struct {
	conf configProvider
}

// SetConfig 关联配置源
func (factory *DefaultBeanFactory) SetConfig(conf interface{}) {
	if c, ok := conf.(configProvider); ok {
		factory.cfg = c
	}
}

// resolveBeanByType 通过具体类型从注册表查找已有Bean，找不到则尝试自动注册并实例化
// 优先级：1) 单例缓存 2) BeanDefinition 3) 自动注册（仅指针型结构体）
func (factory *DefaultBeanFactory) resolveBeanByType(t reflect.Type) interface{} {
	// 通过类型名查找
	beanName := factory.beanNameByType(t)
	if beanName == "" {
		return nil
	}

	if instance := factory.getSingleton(beanName); instance != nil {
		return instance
	}

	if bd, err := factory.registry.GetBeanDefinition(beanName); err == nil {
		instance, err := factory.getBean(beanName, bd)
		if err != nil {
			log.Printf("[hdev-ioc] resolveBeanByType getBean error: %v", err)
			return nil
		}
		return instance
	}

	// 尝试自动注册：仅支持指针型结构体
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		bd := NewBeanDefinition(beanName, t)
		bd.SetScope(ScopeSingleton)
		_ = factory.registry.RegisterBeanDefinition(beanName, bd)
		instance, err := factory.getBean(beanName, bd)
		if err != nil {
			log.Printf("[hdev-ioc] resolveBeanByType auto-register error: %v", err)
			return nil
		}
		return instance
	}

	if t.Kind() == reflect.Struct {
		// 转为指针型尝试自动注册
		ptrType := reflect.PtrTo(t)
		bd := NewBeanDefinition(beanName, ptrType)
		bd.SetScope(ScopeSingleton)
		_ = factory.registry.RegisterBeanDefinition(beanName, bd)
		instance, err := factory.getBean(beanName, bd)
		if err != nil {
			log.Printf("[hdev-ioc] resolveBeanByType auto-register struct error: %v", err)
			return nil
		}
		return instance
	}

	return nil
}

// resolveBeanByInterface 通过接口类型查找已注册Bean中实现该接口的实例
func (factory *DefaultBeanFactory) resolveBeanByInterface(ifaceType reflect.Type) interface{} {
	// 遍历所有已实例化的单例
	factory.singletonMutex.RLock()
	for _, instance := range factory.singletonObjects {
		if instance == nil {
			continue
		}
		if reflect.TypeOf(instance).Implements(ifaceType) {
			factory.singletonMutex.RUnlock()
			return instance
		}
	}
	factory.singletonMutex.RUnlock()

	// 遍历所有已注册Bean定义，对其实例化检查
	for _, beanName := range factory.GetBeanDefinitionNames() {
		bd, err := factory.GetBeanDefinition(beanName)
		if err != nil {
			continue
		}
		if bd.GetBeanType() != nil && bd.GetBeanType().Implements(ifaceType) {
			instance, err := factory.getBean(beanName, bd)
			if err != nil {
				continue
			}
			return instance
		}
	}

	// 自身工厂可作为 IBeanFactory 注入候选
	if reflect.TypeOf(factory).Implements(ifaceType) {
		return factory
	}
	// 配置提供者
	if factory.cfg != nil && reflect.TypeOf(factory.cfg).Implements(ifaceType) {
		return factory.cfg
	}

	return nil
}

// beanNameByType 通过reflect.Type生成BeanName
func (factory *DefaultBeanFactory) beanNameByType(t reflect.Type) string {
	if t == nil {
		return ""
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// injectFields 反射解析struct字段tag，进行依赖注入和值注入
func (factory *DefaultBeanFactory) injectFields(instance interface{}) {
	v := reflect.ValueOf(instance)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		// const tag：直接设置字符串常量
		if cv := fieldType.Tag.Get(TagConst); cv != "" && fieldType.Type.Kind() == reflect.String {
			field.SetString(cv)
			continue
		}

		// value tag：从配置中取值
		if vk := fieldType.Tag.Get(TagValue); vk != "" {
			factory.injectValueTag(field, fieldType, vk)
			continue
		}

		// wired tag：依赖注入
		if wt := fieldType.Tag.Get(TagWired); wt == "true" || wt != "" {
			factory.injectWiredTag(field, fieldType, wt)
			continue
		}
	}
}

// injectValueTag 从配置中取值注入字段
func (factory *DefaultBeanFactory) injectValueTag(field reflect.Value, fieldType reflect.StructField, key string) {
	if factory.cfg == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[hdev-ioc] injectValueTag panic on field %s: %v", fieldType.Name, r)
		}
	}()
	val := factory.cfg.Get(key)
	if val == nil {
		return
	}
	rv := reflect.ValueOf(val)
	if !rv.IsValid() {
		return
	}
	if rv.Type().AssignableTo(field.Type()) {
		field.Set(rv)
	} else if rv.Type().ConvertibleTo(field.Type()) {
		field.Set(rv.Convert(field.Type()))
	}
}

// injectWiredTag 通过类型/接口/特殊预置类型完成依赖注入
// wt：tag值。值为"true"表示按字段类型自动注入；其它值视作Bean名称
func (factory *DefaultBeanFactory) injectWiredTag(field reflect.Value, fieldType reflect.StructField, wt string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[hdev-ioc] injectWiredTag panic on field %s: %v", fieldType.Name, r)
		}
	}()

	ft := fieldType.Type

	// 特殊预置类型：context.Context
	if ft == reflect.TypeOf((*context.Context)(nil)).Elem() {
		if factory.context != nil {
			field.Set(reflect.ValueOf(factory.context))
		}
		return
	}

	// 按名称注入（wt 不是 "true" 或空值时视为名称）
	if wt != "" && wt != "true" {
		if bean := factory.GetBean(wt); bean != nil {
			rv := reflect.ValueOf(bean)
			if rv.Type().AssignableTo(field.Type()) {
				field.Set(rv)
			}
		}
		return
	}

	switch ft.Kind() {
	case reflect.Ptr, reflect.Struct:
		if bean := factory.resolveBeanByType(ft); bean != nil {
			rv := reflect.ValueOf(bean)
			if ft.Kind() == reflect.Struct && rv.Kind() == reflect.Ptr {
				field.Set(rv.Elem())
			} else if rv.Type().AssignableTo(field.Type()) {
				field.Set(rv)
			}
		}
	case reflect.Interface:
		if bean := factory.resolveBeanByInterface(ft); bean != nil {
			rv := reflect.ValueOf(bean)
			if rv.Type().Implements(ft) {
				field.Set(rv)
			}
		}
	}
}

// 兼容性：让hdevcore的BeanPostProcessor接口对未实现也可调用（占位）
var _ = hdevcore.BeanPostProcessor(nil)
