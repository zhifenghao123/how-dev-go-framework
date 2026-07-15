package hdevcore

// ApplicationContext 应用上下文接口
type ApplicationContext interface {
	// GetBean 通过名称获取Bean
	GetBean(name string) interface{}
	
	// GetBeanByType 通过类型获取Bean
	GetBeanByType(beanType interface{}) interface{}
	
	// GetBeansByType 通过类型获取所有Bean
	GetBeansByType(beanType interface{}) []interface{}
	
	// ContainsBean 检查是否包含指定名称的Bean
	ContainsBean(name string) bool
	
	// IsSingleton 检查Bean是否为单例
	IsSingleton(name string) bool
	
	// GetType 获取Bean的类型
	GetType(name string) interface{}
	
	// GetAliases 获取Bean的别名
	GetAliases(name string) []string
	
	// GetEnvironment 获取环境
	GetEnvironment() Environment
	
	// PublishEvent 发布事件
	PublishEvent(event ApplicationEvent)
}

// ConfigurableApplicationContext 可配置的应用上下文接口
type ConfigurableApplicationContext interface {
	ApplicationContext
	
	// SetEnvironment 设置环境
	SetEnvironment(env Environment)
	
	// GetBeanFactory 获取Bean工厂
	GetBeanFactory() ConfigurableListableBeanFactory
	
	// Refresh 刷新应用上下文
	Refresh() error
	
	// Close 关闭应用上下文
	Close() error
	
	// IsActive 检查是否激活
	IsActive() bool
}

// BeanFactory Bean工厂接口
type BeanFactory interface {
	// GetBean 通过名称获取Bean
	GetBean(name string) interface{}
	
	// GetBeanByType 通过类型获取Bean
	GetBeanByType(beanType interface{}) interface{}
	
	// ContainsBean 检查是否包含Bean
	ContainsBean(name string) bool
	
	// IsSingleton 检查是否为单例
	IsSingleton(name string) bool
	
	// IsPrototype 检查是否为原型
	IsPrototype(name string) bool
	
	// GetType 获取Bean类型
	GetType(name string) interface{}
	
	// GetAliases 获取别名
	GetAliases(name string) []string
}

// ListableBeanFactory 可列出的Bean工厂接口
type ListableBeanFactory interface {
	BeanFactory
	
	// GetBeanDefinitionNames 获取所有Bean定义名称
	GetBeanDefinitionNames() []string
	
	// GetBeanDefinitionCount 获取Bean定义数量
	GetBeanDefinitionCount() int
	
	// GetBeansByType 通过类型获取所有Bean
	GetBeansByType(beanType interface{}) map[string]interface{}
}

// ConfigurableListableBeanFactory 可配置可列出的Bean工厂接口
type ConfigurableListableBeanFactory interface {
	ListableBeanFactory
	
	// RegisterSingleton 注册单例Bean
	RegisterSingleton(beanName string, singletonObject interface{})
	
	// RegisterBeanDefinition 注册Bean定义
	RegisterBeanDefinition(beanName string, beanDefinition BeanDefinition)
	
	// GetBeanDefinition 获取Bean定义
	GetBeanDefinition(beanName string) BeanDefinition
	
	// ContainsBeanDefinition 检查是否包含Bean定义
	ContainsBeanDefinition(beanName string) bool
	
	// GetBeanDefinitionNames 获取Bean定义名称
	GetBeanDefinitionNames() []string
	
	// PreInstantiateSingletons 预实例化单例
	PreInstantiateSingletons() error
	
	// AddBeanPostProcessor 添加Bean后置处理器
	AddBeanPostProcessor(processor BeanPostProcessor)
}

// BeanDefinition Bean定义接口
type BeanDefinition interface {
	// GetBeanClassName 获取Bean类名
	GetBeanClassName() string
	
	// GetScope 获取作用域
	GetScope() string
	
	// IsSingleton 是否为单例
	IsSingleton() bool
	
	// IsPrototype 是否为原型
	IsPrototype() bool
	
	// IsLazyInit 是否为懒加载
	IsLazyInit() bool
	
	// GetDependsOn 获取依赖关系
	GetDependsOn() []string
	
	// GetPropertyValues 获取属性值
	GetPropertyValues() map[string]interface{}
}

// ApplicationContextInitializer 应用上下文初始化器接口
type ApplicationContextInitializer func(ctx ApplicationContext) error