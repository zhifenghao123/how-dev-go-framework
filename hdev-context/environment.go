package hdevcontext

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// StandardEnvironment 标准环境实现
// 实现了 hdevcore.Environment 与 hdevcore.ConfigurableEnvironment 接口
type StandardEnvironment struct {
	propertySources []hdevcore.PropertySource
	properties      map[string]string
	activeProfiles  []string
	defaultProfiles []string
	mu              sync.RWMutex
}

// NewStandardEnvironment 创建标准环境
func NewStandardEnvironment() *StandardEnvironment {
	env := &StandardEnvironment{
		propertySources: make([]hdevcore.PropertySource, 0),
		properties:      make(map[string]string),
		activeProfiles:  make([]string, 0),
		defaultProfiles: []string{"default"},
	}
	
	env.addPropertySource(NewSystemPropertySource())
	env.addPropertySource(NewEnvironmentVariablesPropertySource())
	
	return env
}

// GetProperty 获取属性值
func (env *StandardEnvironment) GetProperty(key string) string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	
	// 先查内存覆盖
	if v, ok := env.properties[key]; ok {
		return v
	}
	
	// 从后往前查找，后添加的属性源优先级更高
	for i := len(env.propertySources) - 1; i >= 0; i-- {
		v := env.propertySources[i].GetProperty(key)
		if v != nil {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
			return fmt.Sprintf("%v", v)
		}
	}
	
	return ""
}

// GetPropertyWithDefault 获取属性值，带默认值
func (env *StandardEnvironment) GetPropertyWithDefault(key string, defaultValue string) string {
	value := env.GetProperty(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetRequiredProperty 获取必需的属性值
func (env *StandardEnvironment) GetRequiredProperty(key string) string {
	value := env.GetProperty(key)
	if value == "" {
		panic(fmt.Sprintf("Required property '%s' not found", key))
	}
	return value
}

// GetPropertyAs 获取属性值并转换为指定类型
func (env *StandardEnvironment) GetPropertyAs(key string, targetType interface{}) (interface{}, error) {
	value := env.GetProperty(key)
	if value == "" {
		return nil, fmt.Errorf("property %s not found", key)
	}
	
	switch targetType.(type) {
	case string, *string:
		return value, nil
	case int, *int:
		return strconv.Atoi(value)
	case int64, *int64:
		return strconv.ParseInt(value, 10, 64)
	case float64, *float64:
		return strconv.ParseFloat(value, 64)
	case bool, *bool:
		return strconv.ParseBool(value)
	default:
		return value, nil
	}
}

// ContainsProperty 检查是否包含属性
func (env *StandardEnvironment) ContainsProperty(key string) bool {
	return env.GetProperty(key) != ""
}

// SetProperty 设置属性值
func (env *StandardEnvironment) SetProperty(key, value string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.properties[key] = value
}

// AddPropertySource 添加属性源（公开方法，便于扩展）
func (env *StandardEnvironment) AddPropertySource(source hdevcore.PropertySource) {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.propertySources = append(env.propertySources, source)
}

// GetPropertySources 获取属性源列表
func (env *StandardEnvironment) GetPropertySources() []hdevcore.PropertySource {
	env.mu.RLock()
	defer env.mu.RUnlock()
	return append([]hdevcore.PropertySource{}, env.propertySources...)
}

// Merge 合并环境
func (env *StandardEnvironment) Merge(other hdevcore.Environment) {
	if other == nil {
		return
	}
	
	// 合并激活的profile
	for _, profile := range other.GetActiveProfiles() {
		env.AddActiveProfile(profile)
	}
}

// GetActiveProfiles 获取激活的profile
func (env *StandardEnvironment) GetActiveProfiles() []string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	return append([]string{}, env.activeProfiles...)
}

// SetActiveProfiles 设置激活的profile
func (env *StandardEnvironment) SetActiveProfiles(profiles ...string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.activeProfiles = make([]string, len(profiles))
	copy(env.activeProfiles, profiles)
}

// AddActiveProfile 添加激活的profile
func (env *StandardEnvironment) AddActiveProfile(profile string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	for _, p := range env.activeProfiles {
		if p == profile {
			return
		}
	}
	env.activeProfiles = append(env.activeProfiles, profile)
}

// GetDefaultProfiles 获取默认profile
func (env *StandardEnvironment) GetDefaultProfiles() []string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	return append([]string{}, env.defaultProfiles...)
}

// SetDefaultProfiles 设置默认profile（扩展方法）
func (env *StandardEnvironment) SetDefaultProfiles(profiles ...string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.defaultProfiles = make([]string, len(profiles))
	copy(env.defaultProfiles, profiles)
}

// AcceptsProfiles 检查是否接受指定的profile
func (env *StandardEnvironment) AcceptsProfiles(profiles ...string) bool {
	activeProfiles := env.GetActiveProfiles()
	if len(activeProfiles) == 0 {
		activeProfiles = env.GetDefaultProfiles()
	}
	
	for _, profile := range profiles {
		for _, activeProfile := range activeProfiles {
			if profile == activeProfile {
				return true
			}
		}
	}
	
	return false
}

// addPropertySource 内部添加属性源方法
func (env *StandardEnvironment) addPropertySource(source hdevcore.PropertySource) {
	env.propertySources = append(env.propertySources, source)
}

// ----- PropertySource 实现 -----

// SystemPropertySource 系统属性源
type SystemPropertySource struct {
	name string
}

// NewSystemPropertySource 创建系统属性源
func NewSystemPropertySource() *SystemPropertySource {
	return &SystemPropertySource{name: "systemProperties"}
}

func (sps *SystemPropertySource) GetName() string                    { return sps.name }
func (sps *SystemPropertySource) GetProperty(key string) interface{} { return nil }
func (sps *SystemPropertySource) ContainsProperty(key string) bool   { return false }

// EnvironmentVariablesPropertySource 环境变量属性源
type EnvironmentVariablesPropertySource struct {
	name string
}

// NewEnvironmentVariablesPropertySource 创建环境变量属性源
func NewEnvironmentVariablesPropertySource() *EnvironmentVariablesPropertySource {
	return &EnvironmentVariablesPropertySource{name: "environmentVariables"}
}

func (evps *EnvironmentVariablesPropertySource) GetName() string { return evps.name }

func (evps *EnvironmentVariablesPropertySource) GetProperty(key string) interface{} {
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return nil
}

func (evps *EnvironmentVariablesPropertySource) ContainsProperty(key string) bool {
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	_, ok := os.LookupEnv(envKey)
	return ok
}

// MapPropertySource Map属性源
type MapPropertySource struct {
	name   string
	source map[string]string
}

// NewMapPropertySource 创建Map属性源
func NewMapPropertySource(name string, source map[string]string) *MapPropertySource {
	return &MapPropertySource{name: name, source: source}
}

func (mps *MapPropertySource) GetName() string { return mps.name }

func (mps *MapPropertySource) GetProperty(key string) interface{} {
	if v, ok := mps.source[key]; ok {
		return v
	}
	return nil
}

func (mps *MapPropertySource) ContainsProperty(key string) bool {
	_, ok := mps.source[key]
	return ok
}