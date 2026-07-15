package hdevcontext

import (
	"fmt"
	"reflect"
	"sync"

	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// BeanScope 本模块内使用的Bean作用域类型（与hdev-ioc.BeanScope语义一致但不产生依赖）
type BeanScope string

const (
	// ScopeSingleton 单例作用域
	ScopeSingleton BeanScope = "singleton"
	// ScopePrototype 原型作用域
	ScopePrototype BeanScope = "prototype"
)

// DefaultListableBeanFactory 默认可列出的Bean工厂实现
type DefaultListableBeanFactory struct {
	beanDefinitions    map[string]*BeanDefinition
	singletonObjects   map[string]interface{}
	singletonMutex     sync.RWMutex
	beanPostProcessors []hdevcore.BeanPostProcessor
	aliasMap           map[string]string
}

// BeanDefinition Bean定义（hdev-context层定义，实现hdevcore.BeanDefinition接口）
type BeanDefinition struct {
	name              string
	beanType          reflect.Type
	scope             BeanScope
	instance          interface{}
	lazyInit          bool
	dependsOn         []string
	primary           bool
	initMethodName    string
	destroyMethodName string
	propertyValues    map[string]interface{}
}

// NewDefaultListableBeanFactory 创建默认Bean工厂
func NewDefaultListableBeanFactory() *DefaultListableBeanFactory {
	return &DefaultListableBeanFactory{
		beanDefinitions:    make(map[string]*BeanDefinition),
		singletonObjects:   make(map[string]interface{}),
		beanPostProcessors: make([]hdevcore.BeanPostProcessor, 0),
		aliasMap:           make(map[string]string),
	}
}

// GetBean 通过名称获取Bean
func (factory *DefaultListableBeanFactory) GetBean(name string) interface{} {
	actualName := factory.getActualBeanName(name)
	
	definition, exists := factory.beanDefinitions[actualName]
	if !exists {
		panic(fmt.Sprintf("No bean named '%s' found", name))
	}
	
	switch definition.scope {
	case ScopeSingleton:
		return factory.getSingleton(actualName, definition)
	case ScopePrototype:
		return factory.createBean(definition)
	default:
		return factory.getSingleton(actualName, definition)
	}
}

// GetBeanByType 通过类型获取Bean
func (factory *DefaultListableBeanFactory) GetBeanByType(beanType interface{}) interface{} {
	targetType := reflect.TypeOf(beanType)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	
	var candidate *BeanDefinition
	for _, definition := range factory.beanDefinitions {
		if definition.beanType == nil {
			continue
		}
		bt := definition.beanType
		if bt.Kind() == reflect.Ptr {
			bt = bt.Elem()
		}
		if bt.AssignableTo(targetType) || bt == targetType {
			if definition.primary {
				candidate = definition
				break
			}
			if candidate == nil {
				candidate = definition
			}
		}
	}
	
	if candidate == nil {
		panic(fmt.Sprintf("No bean of type '%s' found", targetType))
	}
	
	return factory.GetBean(candidate.name)
}

// GetBeansByType 通过类型获取所有Bean
func (factory *DefaultListableBeanFactory) GetBeansByType(beanType interface{}) map[string]interface{} {
	targetType := reflect.TypeOf(beanType)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	result := make(map[string]interface{})
	
	for name, definition := range factory.beanDefinitions {
		if definition.beanType == nil {
			continue
		}
		bt := definition.beanType
		if bt.Kind() == reflect.Ptr {
			bt = bt.Elem()
		}
		if bt.AssignableTo(targetType) || bt == targetType {
			result[name] = factory.GetBean(name)
		}
	}
	
	return result
}

// ContainsBean 检查是否包含指定名称的Bean
func (factory *DefaultListableBeanFactory) ContainsBean(name string) bool {
	actualName := factory.getActualBeanName(name)
	_, exists := factory.beanDefinitions[actualName]
	return exists
}

// IsSingleton 检查Bean是否为单例
func (factory *DefaultListableBeanFactory) IsSingleton(name string) bool {
	actualName := factory.getActualBeanName(name)
	definition, exists := factory.beanDefinitions[actualName]
	if !exists {
		return false
	}
	return definition.scope == ScopeSingleton
}

// IsPrototype 检查Bean是否为原型
func (factory *DefaultListableBeanFactory) IsPrototype(name string) bool {
	actualName := factory.getActualBeanName(name)
	definition, exists := factory.beanDefinitions[actualName]
	if !exists {
		return false
	}
	return definition.scope == ScopePrototype
}

// GetType 获取Bean的类型
func (factory *DefaultListableBeanFactory) GetType(name string) interface{} {
	actualName := factory.getActualBeanName(name)
	definition, exists := factory.beanDefinitions[actualName]
	if !exists {
		return nil
	}
	return definition.beanType
}

// GetAliases 获取Bean的别名
func (factory *DefaultListableBeanFactory) GetAliases(name string) []string {
	actualName := factory.getActualBeanName(name)
	aliases := make([]string, 0)
	
	for alias, beanName := range factory.aliasMap {
		if beanName == actualName && alias != actualName {
			aliases = append(aliases, alias)
		}
	}
	
	return aliases
}

// GetBeanDefinitionNames 获取所有Bean定义名称
func (factory *DefaultListableBeanFactory) GetBeanDefinitionNames() []string {
	names := make([]string, 0, len(factory.beanDefinitions))
	for name := range factory.beanDefinitions {
		names = append(names, name)
	}
	return names
}

// GetBeanDefinitionCount 获取Bean定义数量
func (factory *DefaultListableBeanFactory) GetBeanDefinitionCount() int {
	return len(factory.beanDefinitions)
}

// RegisterSingleton 注册单例Bean
func (factory *DefaultListableBeanFactory) RegisterSingleton(name string, instance interface{}) {
	beanType := reflect.TypeOf(instance)
	
	definition := &BeanDefinition{
		name:     name,
		beanType: beanType,
		scope:    ScopeSingleton,
		instance: instance,
		primary:  false,
	}
	
	factory.beanDefinitions[name] = definition
	factory.singletonMutex.Lock()
	factory.singletonObjects[name] = instance
	factory.singletonMutex.Unlock()
}

// RegisterBeanDefinition 注册Bean定义
// 注意：为满足 hdevcore.ConfigurableListableBeanFactory 接口，参数类型为 hdevcore.BeanDefinition
func (factory *DefaultListableBeanFactory) RegisterBeanDefinition(name string, definition hdevcore.BeanDefinition) {
	if bd, ok := definition.(*BeanDefinition); ok {
		factory.beanDefinitions[name] = bd
		return
	}
	// 如果传入的是其它实现，则适配为本地BeanDefinition
	factory.beanDefinitions[name] = &BeanDefinition{
		name:     name,
		scope:    BeanScope(definition.GetScope()),
		lazyInit: definition.IsLazyInit(),
	}
}

// GetBeanDefinition 获取Bean定义
func (factory *DefaultListableBeanFactory) GetBeanDefinition(name string) hdevcore.BeanDefinition {
	bd, exists := factory.beanDefinitions[name]
	if !exists {
		return nil
	}
	return bd
}

// ContainsBeanDefinition 检查是否包含Bean定义
func (factory *DefaultListableBeanFactory) ContainsBeanDefinition(name string) bool {
	_, exists := factory.beanDefinitions[name]
	return exists
}

// AddBeanPostProcessor 添加Bean后置处理器
func (factory *DefaultListableBeanFactory) AddBeanPostProcessor(processor hdevcore.BeanPostProcessor) {
	factory.beanPostProcessors = append(factory.beanPostProcessors, processor)
}

// DestroySingletons 销毁所有单例Bean，调用其销毁方法
func (factory *DefaultListableBeanFactory) DestroySingletons() {
	factory.singletonMutex.Lock()
	defer factory.singletonMutex.Unlock()
	
	for name, instance := range factory.singletonObjects {
		definition, exists := factory.beanDefinitions[name]
		if exists && definition != nil {
			factory.invokeDestroyMethod(instance, definition)
		}
	}
	factory.singletonObjects = make(map[string]interface{})
}

// invokeDestroyMethod 调用销毁方法
func (factory *DefaultListableBeanFactory) invokeDestroyMethod(instance interface{}, definition *BeanDefinition) {
	// 1. 调用配置的销毁方法（优先级最高）
	called := false
	if definition.destroyMethodName != "" {
		method := reflect.ValueOf(instance).MethodByName(definition.destroyMethodName)
		if method.IsValid() {
			method.Call(nil)
			called = true
		}
	}
	// 2. 如果尚未调用过且实现了 DisposableBean 接口，则补调
	if !called {
		if disposable, ok := instance.(hdevcore.DisposableBean); ok {
			_ = disposable.Destroy()
		}
	}
}

// PreInstantiateSingletons 预实例化单例Bean
func (factory *DefaultListableBeanFactory) PreInstantiateSingletons() error {
	factory.singletonMutex.Lock()
	defer factory.singletonMutex.Unlock()
	
	for name, definition := range factory.beanDefinitions {
		if definition.scope == ScopeSingleton && !definition.lazyInit {
			if _, exists := factory.singletonObjects[name]; !exists {
				instance := factory.createBean(definition)
				factory.singletonObjects[name] = instance
			}
		}
	}
	
	return nil
}

// getSingleton 获取单例Bean
func (factory *DefaultListableBeanFactory) getSingleton(name string, definition *BeanDefinition) interface{} {
	factory.singletonMutex.RLock()
	instance, exists := factory.singletonObjects[name]
	factory.singletonMutex.RUnlock()
	
	if exists {
		return instance
	}
	
	factory.singletonMutex.Lock()
	defer factory.singletonMutex.Unlock()
	
	instance, exists = factory.singletonObjects[name]
	if exists {
		return instance
	}
	
	instance = factory.createBean(definition)
	factory.singletonObjects[name] = instance
	
	return instance
}

// createBean 创建Bean实例
func (factory *DefaultListableBeanFactory) createBean(definition *BeanDefinition) interface{} {
	if definition.instance != nil {
		return definition.instance
	}
	
	if definition.beanType == nil {
		return nil
	}
	
	// 根据类型创建实例
	var instance interface{}
	if definition.beanType.Kind() == reflect.Ptr {
		instance = reflect.New(definition.beanType.Elem()).Interface()
	} else {
		instance = reflect.New(definition.beanType).Interface()
	}
	
	// Bean后置处理器：实例化前
	for _, processor := range factory.beanPostProcessors {
		if processed, err := processor.PostProcessBeforeInitialization(instance, definition.name); err == nil && processed != nil {
			instance = processed
		}
	}
	
	// 依赖注入
	factory.injectDependencies(instance, definition)
	
	// 初始化方法
	factory.invokeInitMethod(instance, definition)
	
	// Bean后置处理器：实例化后
	for _, processor := range factory.beanPostProcessors {
		if processed, err := processor.PostProcessAfterInitialization(instance, definition.name); err == nil && processed != nil {
			instance = processed
		}
	}
	
	return instance
}

// injectDependencies 注入依赖
func (factory *DefaultListableBeanFactory) injectDependencies(instance interface{}, definition *BeanDefinition) {
	value := reflect.ValueOf(instance)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return
	}
	
	typeOf := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typeOf.Field(i)
		
		if tag, ok := fieldType.Tag.Lookup("autowired"); ok {
			dependency := factory.tryResolveDependency(tag, field)
			
			if dependency != nil && field.CanSet() {
				depValue := reflect.ValueOf(dependency)
				if depValue.Type().AssignableTo(field.Type()) {
					field.Set(depValue)
				}
			}
		}
	}
}

// tryResolveDependency 尝试解析依赖，失败时返回nil而不是panic
func (factory *DefaultListableBeanFactory) tryResolveDependency(tag string, field reflect.Value) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// 找不到依赖Bean时静默返回nil
			result = nil
		}
	}()
	
	if tag == "" {
		// 按类型注入
		if field.CanAddr() {
			return factory.GetBeanByType(field.Addr().Interface())
		}
		return nil
	}
	// 按名称注入
	return factory.GetBean(tag)
}

// invokeInitMethod 调用初始化方法
func (factory *DefaultListableBeanFactory) invokeInitMethod(instance interface{}, definition *BeanDefinition) {
	if definition.initMethodName == "" {
		// 兼容 InitializingBean 接口
		if initBean, ok := instance.(hdevcore.InitializingBean); ok {
			_ = initBean.AfterPropertiesSet()
		}
		return
	}
	
	method := reflect.ValueOf(instance).MethodByName(definition.initMethodName)
	if method.IsValid() {
		method.Call(nil)
	}
}

// getActualBeanName 获取实际的Bean名称（处理别名）
func (factory *DefaultListableBeanFactory) getActualBeanName(name string) string {
	if actualName, exists := factory.aliasMap[name]; exists {
		return actualName
	}
	return name
}

// ----- BeanDefinition 实现 hdevcore.BeanDefinition 接口 -----

// GetBeanClassName 获取Bean类名
func (bd *BeanDefinition) GetBeanClassName() string {
	if bd.beanType == nil {
		return ""
	}
	return bd.beanType.String()
}

// GetScope 获取作用域
func (bd *BeanDefinition) GetScope() string {
	return string(bd.scope)
}

// IsSingleton 是否为单例
func (bd *BeanDefinition) IsSingleton() bool {
	return bd.scope == ScopeSingleton
}

// IsPrototype 是否为原型
func (bd *BeanDefinition) IsPrototype() bool {
	return bd.scope == ScopePrototype
}

// IsLazyInit 是否为懒加载
func (bd *BeanDefinition) IsLazyInit() bool {
	return bd.lazyInit
}

// GetDependsOn 获取依赖关系
func (bd *BeanDefinition) GetDependsOn() []string {
	return bd.dependsOn
}

// GetPropertyValues 获取属性值
func (bd *BeanDefinition) GetPropertyValues() map[string]interface{} {
	if bd.propertyValues == nil {
		return map[string]interface{}{}
	}
	return bd.propertyValues
}

// ----- 链式构造辅助API -----

// NewBeanDefinition 创建Bean定义
func NewBeanDefinition(name string, beanType reflect.Type) *BeanDefinition {
	return &BeanDefinition{
		name:     name,
		beanType: beanType,
		scope:    ScopeSingleton,
		lazyInit: false,
		primary:  false,
	}
}

// SetScope 设置作用域
func (bd *BeanDefinition) SetScope(scope BeanScope) *BeanDefinition {
	bd.scope = scope
	return bd
}

// SetLazyInit 设置懒加载
func (bd *BeanDefinition) SetLazyInit(lazy bool) *BeanDefinition {
	bd.lazyInit = lazy
	return bd
}

// SetPrimary 设置主要Bean
func (bd *BeanDefinition) SetPrimary(primary bool) *BeanDefinition {
	bd.primary = primary
	return bd
}

// SetDependsOn 设置依赖
func (bd *BeanDefinition) SetDependsOn(dependsOn []string) *BeanDefinition {
	bd.dependsOn = dependsOn
	return bd
}

// SetInitMethodName 设置初始化方法名
func (bd *BeanDefinition) SetInitMethodName(methodName string) *BeanDefinition {
	bd.initMethodName = methodName
	return bd
}

// SetDestroyMethodName 设置销毁方法名
func (bd *BeanDefinition) SetDestroyMethodName(methodName string) *BeanDefinition {
	bd.destroyMethodName = methodName
	return bd
}

// GetName 获取名称
func (bd *BeanDefinition) GetName() string {
	return bd.name
}