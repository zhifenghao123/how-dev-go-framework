package hdevconfig

import (
	"fmt"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// ConfigurableEnvironment 可配置环境实现
type ConfigurableEnvironment struct {
	propertySources    []hdevcore.PropertySource
	activeProfiles    []string
	defaultProfiles   []string
	propertyResolvers []PropertyResolver
	mu                sync.RWMutex
}

// NewConfigurableEnvironment 创建可配置环境
func NewConfigurableEnvironment() *ConfigurableEnvironment {
	env := &ConfigurableEnvironment{
		propertySources:  make([]hdevcore.PropertySource, 0),
		activeProfiles:   make([]string, 0),
		defaultProfiles:  []string{"default"},
	}
	
	// 添加默认属性源
	env.addPropertySource(newSystemPropertySource())
	env.addPropertySource(newEnvironmentVariablePropertySource())
	
	// 添加默认属性解析器
	env.propertyResolvers = []PropertyResolver{
		newDefaultPropertyResolver(),
	}
	
	return env
}

// ContainsProperty 检查是否包含指定属性
func (env *ConfigurableEnvironment) ContainsProperty(key string) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()
	
	for i := len(env.propertySources) - 1; i >= 0; i-- {
		if env.propertySources[i].ContainsProperty(key) {
			return true
		}
	}
	
	return false
}

// GetProperty 获取属性值
func (env *ConfigurableEnvironment) GetProperty(key string) string {
	return env.GetPropertyWithDefault(key, "")
}

// GetPropertyWithDefault 获取属性值，如果不存在返回默认值
// 查找顺序：(1) 激活Profile同名的属性源优先；(2) 其它属性源倒序遍历；(3) 默认值
func (env *ConfigurableEnvironment) GetPropertyWithDefault(key string, defaultValue string) string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	
	activeProfiles := env.activeProfiles
	if len(activeProfiles) == 0 {
		activeProfiles = env.defaultProfiles
	}
	
	// 1. 优先查激活Profile命名的属性源
	for _, profile := range activeProfiles {
		for _, ps := range env.propertySources {
			if ps.GetName() == profile {
				if value := ps.GetProperty(key); value != nil {
					return env.resolvePropertyValue(value)
				}
			}
		}
	}
	
	// 2. 倒序遍历其它属性源
	for i := len(env.propertySources) - 1; i >= 0; i-- {
		// 跳过已经在第1步处理过的Profile源
		isActive := false
		for _, profile := range activeProfiles {
			if env.propertySources[i].GetName() == profile {
				isActive = true
				break
			}
		}
		if isActive {
			continue
		}
		if value := env.propertySources[i].GetProperty(key); value != nil {
			return env.resolvePropertyValue(value)
		}
	}
	
	return defaultValue
}

// GetRequiredProperty 获取必需属性值
func (env *ConfigurableEnvironment) GetRequiredProperty(key string) string {
	value := env.GetProperty(key)
	if value == "" {
		panic(fmt.Sprintf("Required property '%s' not found", key))
	}
	return value
}

// GetPropertyAs 获取属性值并转换为指定类型
func (env *ConfigurableEnvironment) GetPropertyAs(key string, targetType interface{}) (interface{}, error) {
	value := env.GetProperty(key)
	
	// 使用属性解析器进行类型转换
	for _, resolver := range env.propertyResolvers {
		if converted, err := resolver.Resolve(value, targetType); err == nil {
			return converted, nil
		}
	}
	
	return value, fmt.Errorf("unable to convert property '%s' to target type", key)
}

// GetActiveProfiles 获取激活的配置环境
func (env *ConfigurableEnvironment) GetActiveProfiles() []string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	
	if len(env.activeProfiles) == 0 {
		return env.defaultProfiles
	}
	return env.activeProfiles
}

// AcceptsProfiles 检查是否接受指定的配置环境
func (env *ConfigurableEnvironment) AcceptsProfiles(profiles ...string) bool {
	activeProfiles := env.GetActiveProfiles()
	
	for _, profile := range profiles {
		for _, activeProfile := range activeProfiles {
			if profile == activeProfile {
				return true
			}
		}
	}
	
	return false
}

// SetProperty 设置属性值
func (env *ConfigurableEnvironment) SetProperty(key, value string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	
	// 查找第一个可写的属性源
	for _, ps := range env.propertySources {
		if writable, ok := ps.(WritablePropertySource); ok {
			writable.SetProperty(key, value)
			return
		}
	}
	
	// 如果没有可写的属性源，创建一个新的Map属性源
	mapSource := NewMapPropertySource("manualProperties", map[string]interface{}{})
	mapSource.SetProperty(key, value)
	env.addPropertySource(mapSource)
}

// AddActiveProfile 添加激活的配置环境
func (env *ConfigurableEnvironment) AddActiveProfile(profile string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	
	for _, existingProfile := range env.activeProfiles {
		if existingProfile == profile {
			return
		}
	}
	
	env.activeProfiles = append(env.activeProfiles, profile)
}

// SetActiveProfiles 设置激活的配置环境
func (env *ConfigurableEnvironment) SetActiveProfiles(profiles ...string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	
	env.activeProfiles = profiles
}

// Merge 合并另一个环境
func (env *ConfigurableEnvironment) Merge(other hdevcore.Environment) {
	env.mu.Lock()
	defer env.mu.Unlock()
	
	// 合并激活的配置环境
	if otherProfiles := other.GetActiveProfiles(); len(otherProfiles) > 0 {
		env.activeProfiles = otherProfiles
	}
	
	// 合并属性源
	if configurable, ok := other.(*ConfigurableEnvironment); ok {
		for _, ps := range configurable.propertySources {
			env.addPropertySource(ps)
		}
	}
}

// AddPropertySource 添加属性源
func (env *ConfigurableEnvironment) AddPropertySource(propertySource hdevcore.PropertySource) {
	env.mu.Lock()
	defer env.mu.Unlock()
	
	env.addPropertySource(propertySource)
}

// GetPropertySources 获取所有属性源
func (env *ConfigurableEnvironment) GetPropertySources() []hdevcore.PropertySource {
	env.mu.RLock()
	defer env.mu.RUnlock()
	
	return env.propertySources
}

// addPropertySource 内部添加属性源方法
func (env *ConfigurableEnvironment) addPropertySource(propertySource hdevcore.PropertySource) {
	env.propertySources = append(env.propertySources, propertySource)
}

// resolvePropertyValue 解析属性值
func (env *ConfigurableEnvironment) resolvePropertyValue(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

// PropertyResolver 属性解析器接口
type PropertyResolver interface {
	Resolve(value interface{}, targetType interface{}) (interface{}, error)
}

// WritablePropertySource 可写属性源接口
type WritablePropertySource interface {
	hdevcore.PropertySource
	SetProperty(key string, value interface{})
}

// 私有实现类

type systemPropertySource struct {
	name string
}

func newSystemPropertySource() *systemPropertySource {
	return &systemPropertySource{
		name: "systemProperties",
	}
}

func (sps *systemPropertySource) GetName() string {
	return sps.name
}

func (sps *systemPropertySource) GetProperty(key string) interface{} {
	// 简化实现
	return nil
}

func (sps *systemPropertySource) ContainsProperty(key string) bool {
	return false
}

type environmentVariablePropertySource struct {
	name string
}

func newEnvironmentVariablePropertySource() *environmentVariablePropertySource {
	return &environmentVariablePropertySource{
		name: "environmentVariables",
	}
}

func (evps *environmentVariablePropertySource) GetName() string {
	return evps.name
}

func (evps *environmentVariablePropertySource) GetProperty(key string) interface{} {
	// 简化实现
	return nil
}

func (evps *environmentVariablePropertySource) ContainsProperty(key string) bool {
	return false
}

type defaultPropertyResolver struct{}

func newDefaultPropertyResolver() *defaultPropertyResolver {
	return &defaultPropertyResolver{}
}

func (dpr *defaultPropertyResolver) Resolve(value interface{}, targetType interface{}) (interface{}, error) {
	// 简化实现
	return value, nil
}