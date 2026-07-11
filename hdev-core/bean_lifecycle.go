package hdevcore

// BeanNameAware Bean名称感知接口
type BeanNameAware interface {
	// SetBeanName 设置Bean名称
	SetBeanName(name string)
}

// BeanFactoryAware Bean工厂感知接口
type BeanFactoryAware interface {
	// SetBeanFactory 设置Bean工厂
	SetBeanFactory(beanFactory BeanFactory)
}

// ApplicationContextAware 应用上下文感知接口
type ApplicationContextAware interface {
	// SetApplicationContext 设置应用上下文
	SetApplicationContext(ctx ApplicationContext)
}

// EnvironmentAware 环境感知接口
type EnvironmentAware interface {
	// SetEnvironment 设置环境
	SetEnvironment(env Environment)
}

// InitializingBean 初始化Bean接口
type InitializingBean interface {
	// AfterPropertiesSet 属性设置完成后调用
	AfterPropertiesSet() error
}

// DisposableBean 销毁Bean接口
type DisposableBean interface {
	// Destroy Bean销毁时调用
	Destroy() error
}

// BeanPostProcessor Bean后置处理器接口
type BeanPostProcessor interface {
	// PostProcessBeforeInitialization 初始化前处理
	PostProcessBeforeInitialization(bean interface{}, beanName string) (interface{}, error)
	
	// PostProcessAfterInitialization 初始化后处理
	PostProcessAfterInitialization(bean interface{}, beanName string) (interface{}, error)
}

// BeanFactoryPostProcessor Bean工厂后置处理器接口
type BeanFactoryPostProcessor interface {
	// PostProcessBeanFactory Bean工厂后处理
	PostProcessBeanFactory(beanFactory ConfigurableListableBeanFactory) error
}

// Ordered 排序接口
type Ordered interface {
	// GetOrder 获取排序值
	GetOrder() int
}

// PriorityOrdered 优先级排序接口
type PriorityOrdered interface {
	Ordered
}