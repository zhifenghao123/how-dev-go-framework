package hdev_ioc

import (
	"fmt"
	"reflect"
	"sync"
)

// BeanScope Bean作用域
type BeanScope string

const (
	// ScopeSingleton 单例作用域
	ScopeSingleton BeanScope = "singleton"
	// ScopePrototype 原型作用域
	ScopePrototype BeanScope = "prototype"
)

// BeanDefinition Bean定义
type BeanDefinition struct {
	name           string
	beanType       reflect.Type
	instance       interface{}
	scope          BeanScope
	lazyInit       bool
	dependsOn      []string
	propertyValues map[string]interface{}
	factoryMethod  string
	factoryBean    string
	initMethod     string
	destroyMethod  string
	primary        bool
	mu             sync.RWMutex
}

// NewBeanDefinition 创建Bean定义
func NewBeanDefinition(name string, beanType reflect.Type) *BeanDefinition {
	return &BeanDefinition{
		name:           name,
		beanType:       beanType,
		scope:          ScopeSingleton,
		lazyInit:       false,
		dependsOn:      make([]string, 0),
		propertyValues: make(map[string]interface{}),
		primary:        false,
	}
}

// GetName 获取Bean名称
func (bd *BeanDefinition) GetName() string {
	return bd.name
}

// GetBeanType 获取Bean类型
func (bd *BeanDefinition) GetBeanType() reflect.Type {
	return bd.beanType
}

// GetInstance 获取Bean实例
func (bd *BeanDefinition) GetInstance() interface{} {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	
	return bd.instance
}

// SetInstance 设置Bean实例
func (bd *BeanDefinition) SetInstance(instance interface{}) {
	bd.mu.Lock()
	defer bd.mu.Unlock()
	
	bd.instance = instance
}

// GetScope 获取Bean作用域
func (bd *BeanDefinition) GetScope() BeanScope {
	return bd.scope
}

// SetScope 设置Bean作用域
func (bd *BeanDefinition) SetScope(scope BeanScope) {
	bd.scope = scope
}

// IsLazyInit 是否延迟初始化
func (bd *BeanDefinition) IsLazyInit() bool {
	return bd.lazyInit
}

// SetLazyInit 设置延迟初始化
func (bd *BeanDefinition) SetLazyInit(lazyInit bool) {
	bd.lazyInit = lazyInit
}

// GetDependsOn 获取依赖的Bean名称
func (bd *BeanDefinition) GetDependsOn() []string {
	return bd.dependsOn
}

// SetDependsOn 设置依赖的Bean名称
func (bd *BeanDefinition) SetDependsOn(dependsOn []string) {
	bd.dependsOn = dependsOn
}

// AddDependsOn 添加依赖的Bean名称
func (bd *BeanDefinition) AddDependsOn(beanName string) {
	bd.dependsOn = append(bd.dependsOn, beanName)
}

// GetPropertyValues 获取属性值
func (bd *BeanDefinition) GetPropertyValues() map[string]interface{} {
	return bd.propertyValues
}

// SetPropertyValue 设置属性值
func (bd *BeanDefinition) SetPropertyValue(name string, value interface{}) {
	bd.propertyValues[name] = value
}

// GetPropertyValue 获取属性值
func (bd *BeanDefinition) GetPropertyValue(name string) interface{} {
	return bd.propertyValues[name]
}

// GetFactoryMethod 获取工厂方法
func (bd *BeanDefinition) GetFactoryMethod() string {
	return bd.factoryMethod
}

// SetFactoryMethod 设置工厂方法
func (bd *BeanDefinition) SetFactoryMethod(method string) {
	bd.factoryMethod = method
}

// GetFactoryBean 获取工厂Bean
func (bd *BeanDefinition) GetFactoryBean() string {
	return bd.factoryBean
}

// SetFactoryBean 设置工厂Bean
func (bd *BeanDefinition) SetFactoryBean(beanName string) {
	bd.factoryBean = beanName
}

// GetInitMethod 获取初始化方法
func (bd *BeanDefinition) GetInitMethod() string {
	return bd.initMethod
}

// SetInitMethod 设置初始化方法
func (bd *BeanDefinition) SetInitMethod(method string) {
	bd.initMethod = method
}

// GetDestroyMethod 获取销毁方法
func (bd *BeanDefinition) GetDestroyMethod() string {
	return bd.destroyMethod
}

// SetDestroyMethod 设置销毁方法
func (bd *BeanDefinition) SetDestroyMethod(method string) {
	bd.destroyMethod = method
}

// SetInitMethodName 设置初始化方法名（Spring风格别名）
func (bd *BeanDefinition) SetInitMethodName(method string) {
	bd.initMethod = method
}

// SetDestroyMethodName 设置销毁方法名（Spring风格别名）
func (bd *BeanDefinition) SetDestroyMethodName(method string) {
	bd.destroyMethod = method
}

// IsPrimary 是否为主要Bean
func (bd *BeanDefinition) IsPrimary() bool {
	return bd.primary
}

// SetPrimary 设置为主要Bean
func (bd *BeanDefinition) SetPrimary(primary bool) {
	bd.primary = primary
}

// BeanDefinitionBuilder Bean定义构建器
type BeanDefinitionBuilder struct {
	definition *BeanDefinition
}

// NewBeanDefinitionBuilder 创建Bean定义构建器
func NewBeanDefinitionBuilder(name string, beanType reflect.Type) *BeanDefinitionBuilder {
	return &BeanDefinitionBuilder{
		definition: NewBeanDefinition(name, beanType),
	}
}

// WithScope 设置作用域
func (builder *BeanDefinitionBuilder) WithScope(scope BeanScope) *BeanDefinitionBuilder {
	builder.definition.SetScope(scope)
	return builder
}

// WithLazyInit 设置延迟初始化
func (builder *BeanDefinitionBuilder) WithLazyInit(lazyInit bool) *BeanDefinitionBuilder {
	builder.definition.SetLazyInit(lazyInit)
	return builder
}

// WithDependsOn 设置依赖
func (builder *BeanDefinitionBuilder) WithDependsOn(dependsOn ...string) *BeanDefinitionBuilder {
	builder.definition.SetDependsOn(dependsOn)
	return builder
}

// WithPropertyValue 设置属性值
func (builder *BeanDefinitionBuilder) WithPropertyValue(name string, value interface{}) *BeanDefinitionBuilder {
	builder.definition.SetPropertyValue(name, value)
	return builder
}

// WithFactoryMethod 设置工厂方法
func (builder *BeanDefinitionBuilder) WithFactoryMethod(method string) *BeanDefinitionBuilder {
	builder.definition.SetFactoryMethod(method)
	return builder
}

// WithFactoryBean 设置工厂Bean
func (builder *BeanDefinitionBuilder) WithFactoryBean(beanName string) *BeanDefinitionBuilder {
	builder.definition.SetFactoryBean(beanName)
	return builder
}

// WithInitMethod 设置初始化方法
func (builder *BeanDefinitionBuilder) WithInitMethod(method string) *BeanDefinitionBuilder {
	builder.definition.SetInitMethod(method)
	return builder
}

// WithDestroyMethod 设置销毁方法
func (builder *BeanDefinitionBuilder) WithDestroyMethod(method string) *BeanDefinitionBuilder {
	builder.definition.SetDestroyMethod(method)
	return builder
}

// WithPrimary 设置为主要Bean
func (builder *BeanDefinitionBuilder) WithPrimary(primary bool) *BeanDefinitionBuilder {
	builder.definition.SetPrimary(primary)
	return builder
}

// Build 构建Bean定义
func (builder *BeanDefinitionBuilder) Build() *BeanDefinition {
	return builder.definition
}

// BeanDefinitionRegistry Bean定义注册表接口
type BeanDefinitionRegistry interface {
	// RegisterBeanDefinition 注册Bean定义
	RegisterBeanDefinition(beanName string, beanDefinition *BeanDefinition) error
	// RemoveBeanDefinition 移除Bean定义
	RemoveBeanDefinition(beanName string) error
	// GetBeanDefinition 获取Bean定义
	GetBeanDefinition(beanName string) (*BeanDefinition, error)
	// ContainsBeanDefinition 检查是否包含Bean定义
	ContainsBeanDefinition(beanName string) bool
	// GetBeanDefinitionNames 获取所有Bean定义名称
	GetBeanDefinitionNames() []string
	// GetBeanDefinitionCount 获取Bean定义数量
	GetBeanDefinitionCount() int
}

// DefaultBeanDefinitionRegistry 默认Bean定义注册表
type DefaultBeanDefinitionRegistry struct {
	beanDefinitions map[string]*BeanDefinition
	mu              sync.RWMutex
}

// NewDefaultBeanDefinitionRegistry 创建默认Bean定义注册表
func NewDefaultBeanDefinitionRegistry() *DefaultBeanDefinitionRegistry {
	return &DefaultBeanDefinitionRegistry{
		beanDefinitions: make(map[string]*BeanDefinition),
	}
}

// RegisterBeanDefinition 注册Bean定义
func (registry *DefaultBeanDefinitionRegistry) RegisterBeanDefinition(beanName string, beanDefinition *BeanDefinition) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	
	if _, exists := registry.beanDefinitions[beanName]; exists {
		return fmt.Errorf("Bean definition already exists: " + beanName)
	}
	
	registry.beanDefinitions[beanName] = beanDefinition
	return nil
}

// RemoveBeanDefinition 移除Bean定义
func (registry *DefaultBeanDefinitionRegistry) RemoveBeanDefinition(beanName string) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	
	if _, exists := registry.beanDefinitions[beanName]; !exists {
		return fmt.Errorf("Bean definition not found: " + beanName)
	}
	
	delete(registry.beanDefinitions, beanName)
	return nil
}

// GetBeanDefinition 获取Bean定义
func (registry *DefaultBeanDefinitionRegistry) GetBeanDefinition(beanName string) (*BeanDefinition, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	
	beanDefinition, exists := registry.beanDefinitions[beanName]
	if !exists {
		return nil, fmt.Errorf("Bean definition not found: " + beanName)
	}
	
	return beanDefinition, nil
}

// ContainsBeanDefinition 检查是否包含Bean定义
func (registry *DefaultBeanDefinitionRegistry) ContainsBeanDefinition(beanName string) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	
	_, exists := registry.beanDefinitions[beanName]
	return exists
}

// GetBeanDefinitionNames 获取所有Bean定义名称
func (registry *DefaultBeanDefinitionRegistry) GetBeanDefinitionNames() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	
	names := make([]string, 0, len(registry.beanDefinitions))
	for name := range registry.beanDefinitions {
		names = append(names, name)
	}
	return names
}

// GetBeanDefinitionCount 获取Bean定义数量
func (registry *DefaultBeanDefinitionRegistry) GetBeanDefinitionCount() int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	
	return len(registry.beanDefinitions)
}

// BeanDefinitionException Bean定义异常
type BeanDefinitionException struct {
	message string
}

// NewBeanDefinitionException 创建Bean定义异常
func NewBeanDefinitionException(message string) *BeanDefinitionException {
	return &BeanDefinitionException{message: message}
}

// Error 实现error接口
func (e *BeanDefinitionException) Error() string {
	return e.message
}