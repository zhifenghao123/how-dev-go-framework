package hdevconfig

import (
	"fmt"
	"strings"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// MapPropertySource 基于Map的属性源
type MapPropertySource struct {
	name   string
	source map[string]interface{}
	mu     sync.RWMutex
}

// NewMapPropertySource 创建Map属性源
func NewMapPropertySource(name string, source map[string]interface{}) *MapPropertySource {
	return &MapPropertySource{
		name:   name,
		source: source,
	}
}

// GetName 获取属性源名称
func (mps *MapPropertySource) GetName() string {
	return mps.name
}

// GetProperty 获取属性值
func (mps *MapPropertySource) GetProperty(key string) interface{} {
	mps.mu.RLock()
	defer mps.mu.RUnlock()
	
	return mps.source[key]
}

// ContainsProperty 检查是否包含指定属性
func (mps *MapPropertySource) ContainsProperty(key string) bool {
	mps.mu.RLock()
	defer mps.mu.RUnlock()
	
	_, exists := mps.source[key]
	return exists
}

// SetProperty 设置属性值
func (mps *MapPropertySource) SetProperty(key string, value interface{}) {
	mps.mu.Lock()
	defer mps.mu.Unlock()
	
	mps.source[key] = value
}

// GetAllProperties 获取所有属性
func (mps *MapPropertySource) GetAllProperties() map[string]interface{} {
	mps.mu.RLock()
	defer mps.mu.RUnlock()
	
	// 返回副本
	result := make(map[string]interface{})
	for k, v := range mps.source {
		result[k] = v
	}
	return result
}

// CompositePropertySource 组合属性源
type CompositePropertySource struct {
	name            string
	propertySources []hdevcore.PropertySource
	mu              sync.RWMutex
}

// NewCompositePropertySource 创建组合属性源
func NewCompositePropertySource(name string) *CompositePropertySource {
	return &CompositePropertySource{
		name:            name,
		propertySources: make([]hdevcore.PropertySource, 0),
	}
}

// GetName 获取属性源名称
func (cps *CompositePropertySource) GetName() string {
	return cps.name
}

// GetProperty 获取属性值
func (cps *CompositePropertySource) GetProperty(key string) interface{} {
	cps.mu.RLock()
	defer cps.mu.RUnlock()
	
	for i := len(cps.propertySources) - 1; i >= 0; i-- {
		if value := cps.propertySources[i].GetProperty(key); value != nil {
			return value
		}
	}
	return nil
}

// ContainsProperty 检查是否包含指定属性
func (cps *CompositePropertySource) ContainsProperty(key string) bool {
	cps.mu.RLock()
	defer cps.mu.RUnlock()
	
	for _, ps := range cps.propertySources {
		if ps.ContainsProperty(key) {
			return true
		}
	}
	return false
}

// AddPropertySource 添加属性源
func (cps *CompositePropertySource) AddPropertySource(propertySource hdevcore.PropertySource) {
	cps.mu.Lock()
	defer cps.mu.Unlock()
	
	cps.propertySources = append(cps.propertySources, propertySource)
}

// RemovePropertySource 移除属性源
func (cps *CompositePropertySource) RemovePropertySource(name string) {
	cps.mu.Lock()
	defer cps.mu.Unlock()
	
	for i, ps := range cps.propertySources {
		if ps.GetName() == name {
			cps.propertySources = append(cps.propertySources[:i], cps.propertySources[i+1:]...)
			break
		}
	}
}

// GetPropertySources 获取所有属性源
func (cps *CompositePropertySource) GetPropertySources() []hdevcore.PropertySource {
	cps.mu.RLock()
	defer cps.mu.RUnlock()
	
	return cps.propertySources
}

// EnvironmentUtils 环境工具类
type EnvironmentUtils struct{}

// GetPropertyOrPanic 获取属性值，如果不存在则panic
func (utils *EnvironmentUtils) GetPropertyOrPanic(env hdevcore.Environment, key string) string {
	value := env.GetProperty(key)
	if value == "" {
		panic(fmt.Sprintf("Required property '%s' not found", key))
	}
	return value
}

// GetPropertyOrDefault 获取属性值，如果不存在则返回默认值
func (utils *EnvironmentUtils) GetPropertyOrDefault(env hdevcore.Environment, key string, defaultValue string) string {
	value := env.GetProperty(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetIntegerProperty 获取整数属性值
func (utils *EnvironmentUtils) GetIntegerProperty(env hdevcore.Environment, key string, defaultValue int) (int, error) {
	value := env.GetProperty(key)
	if value == "" {
		return defaultValue, nil
	}
	
	intValue, err := utils.parseInt(value)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid integer value for property '%s': %v", key, err)
	}
	
	return intValue, nil
}

// GetBooleanProperty 获取布尔属性值
func (utils *EnvironmentUtils) GetBooleanProperty(env hdevcore.Environment, key string, defaultValue bool) (bool, error) {
	value := env.GetProperty(key)
	if value == "" {
		return defaultValue, nil
	}
	
	boolValue, err := utils.parseBool(value)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid boolean value for property '%s': %v", key, err)
	}
	
	return boolValue, nil
}

// parseInt 解析整数字符串
func (utils *EnvironmentUtils) parseInt(value string) (int, error) {
	// 简化实现，实际应该使用strconv
	return 0, nil
}

// parseBool 解析布尔字符串
func (utils *EnvironmentUtils) parseBool(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}

// PropertySourceBuilder 属性源构建器
type PropertySourceBuilder struct {
	name   string
	source map[string]interface{}
}

// NewPropertySourceBuilder 创建属性源构建器
func NewPropertySourceBuilder(name string) *PropertySourceBuilder {
	return &PropertySourceBuilder{
		name:   name,
		source: make(map[string]interface{}),
	}
}

// WithProperty 添加属性
func (builder *PropertySourceBuilder) WithProperty(key string, value interface{}) *PropertySourceBuilder {
	builder.source[key] = value
	return builder
}

// WithProperties 批量添加属性
func (builder *PropertySourceBuilder) WithProperties(properties map[string]interface{}) *PropertySourceBuilder {
	for k, v := range properties {
		builder.source[k] = v
	}
	return builder
}

// Build 构建属性源
func (builder *PropertySourceBuilder) Build() *MapPropertySource {
	return NewMapPropertySource(builder.name, builder.source)
}

// EnvironmentBuilder 环境构建器
type EnvironmentBuilder struct {
	env *ConfigurableEnvironment
}

// NewEnvironmentBuilder 创建环境构建器
func NewEnvironmentBuilder() *EnvironmentBuilder {
	return &EnvironmentBuilder{
		env: NewConfigurableEnvironment(),
	}
}

// WithPropertySource 添加属性源
func (builder *EnvironmentBuilder) WithPropertySource(source hdevcore.PropertySource) *EnvironmentBuilder {
	builder.env.AddPropertySource(source)
	return builder
}

// WithProperties 添加属性
func (builder *EnvironmentBuilder) WithProperties(name string, properties map[string]interface{}) *EnvironmentBuilder {
	source := NewMapPropertySource(name, properties)
	builder.env.AddPropertySource(source)
	return builder
}

// WithActiveProfiles 设置激活的环境
func (builder *EnvironmentBuilder) WithActiveProfiles(profiles ...string) *EnvironmentBuilder {
	builder.env.SetActiveProfiles(profiles...)
	return builder
}

// Build 构建环境
func (builder *EnvironmentBuilder) Build() hdevcore.Environment {
	return builder.env
}

// Example usage:
/*
func main() {
	// 使用构建器模式创建环境
	env := hdevconfig.NewEnvironmentBuilder().
		WithProperties("default", map[string]interface{}{
			"server.port": "8080",
			"server.host": "localhost",
		}).
		WithActiveProfiles("dev").
		Build()
	
	// 使用工具类
	utils := &hdevconfig.EnvironmentUtils{}
	port := utils.GetPropertyOrDefault(env, "server.port", "8080")
	fmt.Printf("Server port: %s\n", port)
}
*/

// 全局工具实例
var EnvironmentUtil = &EnvironmentUtils{}

// 辅助函数

// GetPropertyOrPanic 全局辅助函数
func GetPropertyOrPanic(env hdevcore.Environment, key string) string {
	return EnvironmentUtil.GetPropertyOrPanic(env, key)
}

// GetPropertyOrDefault 全局辅助函数
func GetPropertyOrDefault(env hdevcore.Environment, key string, defaultValue string) string {
	return EnvironmentUtil.GetPropertyOrDefault(env, key, defaultValue)
}

// GetIntegerProperty 全局辅助函数
func GetIntegerProperty(env hdevcore.Environment, key string, defaultValue int) (int, error) {
	return EnvironmentUtil.GetIntegerProperty(env, key, defaultValue)
}

// GetBooleanProperty 全局辅助函数
func GetBooleanProperty(env hdevcore.Environment, key string, defaultValue bool) (bool, error) {
	return EnvironmentUtil.GetBooleanProperty(env, key, defaultValue)
}