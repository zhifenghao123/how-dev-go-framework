package hdevcore

// Environment 环境接口，提供对应用环境的访问
// 参考Spring的Environment接口设计
type Environment interface {
	// ContainsProperty 检查是否包含指定属性
	ContainsProperty(key string) bool

	// GetProperty 获取属性值，如果不存在返回空字符串
	GetProperty(key string) string

	// GetPropertyWithDefault 获取属性值，如果不存在返回默认值
	GetPropertyWithDefault(key string, defaultValue string) string

	// GetRequiredProperty 获取必需属性值，如果不存在则panic
	GetRequiredProperty(key string) string

	// GetPropertyAs 获取属性值并转换为指定类型
	GetPropertyAs(key string, targetType interface{}) (interface{}, error)

	// GetActiveProfiles 获取激活的配置环境
	GetActiveProfiles() []string

	// AcceptsProfiles 检查是否接受指定的配置环境
	AcceptsProfiles(profiles ...string) bool
}

// ConfigurableEnvironment 可配置的环境接口
type ConfigurableEnvironment interface {
	Environment

	// SetProperty 设置属性值
	SetProperty(key, value string)

	// AddActiveProfile 添加激活的配置环境
	AddActiveProfile(profile string)

	// SetActiveProfiles 设置激活的配置环境
	SetActiveProfiles(profiles ...string)

	// Merge 合并另一个环境
	Merge(other Environment)
}

// PropertySource 属性源接口
type PropertySource interface {
	// GetName 获取属性源名称
	GetName() string

	// GetProperty 获取属性值
	GetProperty(key string) interface{}

	// ContainsProperty 检查是否包含属性
	ContainsProperty(key string) bool
}

// PropertyResolver 属性解析器接口
type PropertyResolver interface {
	// ResolvePlaceholders 解析占位符
	ResolvePlaceholders(text string) (string, error)

	// ResolveRequiredPlaceholders 解析必需的占位符
	ResolveRequiredPlaceholders(text string) (string, error)
}
