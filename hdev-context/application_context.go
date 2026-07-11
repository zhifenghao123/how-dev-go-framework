package hdevcontext

import (
	"sync"
	"time"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// DefaultApplicationContext 默认应用上下文实现
type DefaultApplicationContext struct {
	parent              hdevcore.ApplicationContext
	beanFactory         hdevcore.ConfigurableListableBeanFactory
	environment         hdevcore.ConfigurableEnvironment
	eventMulticaster    hdevcore.ApplicationEventMulticaster
	startupDate         int64
	active             bool
	closed             bool
	mu                 sync.RWMutex
	listenerBeans      []string
	earlyApplicationEvents []hdevcore.ApplicationEvent
}

// NewDefaultApplicationContext 创建默认应用上下文
func NewDefaultApplicationContext() *DefaultApplicationContext {
	ctx := &DefaultApplicationContext{
		startupDate: time.Now().UnixNano(),
		active:      false,
		closed:      true,
	}
	
	// 初始化Bean工厂
	ctx.beanFactory = NewDefaultListableBeanFactory()
	
	// 初始化环境
	ctx.environment = NewStandardEnvironment()
	
	// 初始化事件广播器
	ctx.eventMulticaster = NewSimpleApplicationEventMulticaster()
	
	return ctx
}

// GetBean 通过名称获取Bean
func (ctx *DefaultApplicationContext) GetBean(name string) interface{} {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		panic("ApplicationContext is not active")
	}
	
	return ctx.beanFactory.GetBean(name)
}

// GetBeanByType 通过类型获取Bean
func (ctx *DefaultApplicationContext) GetBeanByType(beanType interface{}) interface{} {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		panic("ApplicationContext is not active")
	}
	
	return ctx.beanFactory.GetBeanByType(beanType)
}

// GetBeansByType 通过类型获取所有Bean
func (ctx *DefaultApplicationContext) GetBeansByType(beanType interface{}) []interface{} {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		panic("ApplicationContext is not active")
	}
	
	beansMap := ctx.beanFactory.GetBeansByType(beanType)
	beans := make([]interface{}, 0, len(beansMap))
	for _, bean := range beansMap {
		beans = append(beans, bean)
	}
	return beans
}

// ContainsBean 检查是否包含指定名称的Bean
func (ctx *DefaultApplicationContext) ContainsBean(name string) bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		return false
	}
	
	return ctx.beanFactory.ContainsBean(name)
}

// IsSingleton 检查Bean是否为单例
func (ctx *DefaultApplicationContext) IsSingleton(name string) bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		return false
	}
	
	return ctx.beanFactory.IsSingleton(name)
}

// GetType 获取Bean的类型
func (ctx *DefaultApplicationContext) GetType(name string) interface{} {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		return nil
	}
	
	return ctx.beanFactory.GetType(name)
}

// GetAliases 获取Bean的别名
func (ctx *DefaultApplicationContext) GetAliases(name string) []string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		return nil
	}
	
	return ctx.beanFactory.GetAliases(name)
}

// GetEnvironment 获取环境
func (ctx *DefaultApplicationContext) GetEnvironment() hdevcore.Environment {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	return ctx.environment
}

// PublishEvent 发布事件
func (ctx *DefaultApplicationContext) PublishEvent(event hdevcore.ApplicationEvent) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	if !ctx.active {
		// 如果上下文还未激活，缓存早期事件
		ctx.earlyApplicationEvents = append(ctx.earlyApplicationEvents, event)
		return
	}
	
	ctx.eventMulticaster.PublishEvent(event)
}

// SetEnvironment 设置环境
func (ctx *DefaultApplicationContext) SetEnvironment(env hdevcore.Environment) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	
	if configurableEnv, ok := env.(hdevcore.ConfigurableEnvironment); ok {
		ctx.environment = configurableEnv
	} else {
		// 如果不是可配置环境，创建包装器
		ctx.environment = NewStandardEnvironment()
		ctx.environment.Merge(env)
	}
}

// GetBeanFactory 获取Bean工厂
func (ctx *DefaultApplicationContext) GetBeanFactory() hdevcore.ConfigurableListableBeanFactory {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	return ctx.beanFactory
}

// Refresh 刷新应用上下文
func (ctx *DefaultApplicationContext) Refresh() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	
	if ctx.active {
		// 如果已经激活，先关闭
		ctx.doClose()
	}
	
	// 准备刷新
	err := ctx.prepareRefresh()
	if err != nil {
		return err
	}
	
	// 配置Bean工厂
	err = ctx.prepareBeanFactory(ctx.beanFactory)
	if err != nil {
		return err
	}
	
	// 调用Bean工厂后置处理器
	err = ctx.invokeBeanFactoryPostProcessors(ctx.beanFactory)
	if err != nil {
		return err
	}
	
	// 注册Bean后置处理器
	err = ctx.registerBeanPostProcessors(ctx.beanFactory)
	if err != nil {
		return err
	}
	
	// 初始化消息源
	err = ctx.initMessageSource()
	if err != nil {
		return err
	}
	
	// 初始化事件广播器
	err = ctx.initApplicationEventMulticaster()
	if err != nil {
		return err
	}
	
	// 注册监听器
	err = ctx.registerListeners()
	if err != nil {
		return err
	}
	
	// 实例化剩余的单例Bean
	err = ctx.finishBeanFactoryInitialization(ctx.beanFactory)
	if err != nil {
		return err
	}
	
	// 完成刷新
	err = ctx.finishRefresh()
	if err != nil {
		return err
	}
	
	ctx.active = true
	ctx.closed = false
	
	return nil
}

// Close 关闭应用上下文
func (ctx *DefaultApplicationContext) Close() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	
	return ctx.doClose()
}

// IsActive 检查是否激活
func (ctx *DefaultApplicationContext) IsActive() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	
	return ctx.active
}

// doClose 执行关闭操作
func (ctx *DefaultApplicationContext) doClose() error {
	if !ctx.active {
		return nil
	}
	
	// 发布上下文关闭事件
	ctx.publishEvent(newContextClosedEvent(ctx))
	
	// 销毁单例Bean
	ctx.destroyBeans()
	
	// 关闭Bean工厂
	ctx.beanFactory = nil
	
	ctx.active = false
	ctx.closed = true
	
	return nil
}

// prepareRefresh 准备刷新
func (ctx *DefaultApplicationContext) prepareRefresh() error {
	// 初始化启动时间
	ctx.startupDate = time.Now().UnixNano()
	
	// 初始化早期事件
	ctx.earlyApplicationEvents = make([]hdevcore.ApplicationEvent, 0)
	
	// 验证环境
	if ctx.environment == nil {
		ctx.environment = NewStandardEnvironment()
	}
	
	return nil
}

// prepareBeanFactory 准备Bean工厂
func (ctx *DefaultApplicationContext) prepareBeanFactory(beanFactory hdevcore.ConfigurableListableBeanFactory) error {
	// 注册环境Bean
	beanFactory.RegisterSingleton("environment", ctx.environment)
	
	// 注册应用上下文Bean
	beanFactory.RegisterSingleton("applicationContext", ctx)
	
	return nil
}

// invokeBeanFactoryPostProcessors 调用Bean工厂后置处理器
func (ctx *DefaultApplicationContext) invokeBeanFactoryPostProcessors(beanFactory hdevcore.ConfigurableListableBeanFactory) error {
	// 直接从工厂取，避免反rlock导致死锁
	postProcessors := beanFactory.GetBeansByType((*hdevcore.BeanFactoryPostProcessor)(nil))
	
	for _, postProcessor := range postProcessors {
		if factoryPostProcessor, ok := postProcessor.(hdevcore.BeanFactoryPostProcessor); ok {
			if err := factoryPostProcessor.PostProcessBeanFactory(beanFactory); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// registerBeanPostProcessors 注册Bean后置处理器
func (ctx *DefaultApplicationContext) registerBeanPostProcessors(beanFactory hdevcore.ConfigurableListableBeanFactory) error {
	// 直接从工厂取，避免反rlock导致死锁
	postProcessors := beanFactory.GetBeansByType((*hdevcore.BeanPostProcessor)(nil))
	
	for _, postProcessor := range postProcessors {
		if beanPostProcessor, ok := postProcessor.(hdevcore.BeanPostProcessor); ok {
			beanFactory.AddBeanPostProcessor(beanPostProcessor)
		}
	}
	
	return nil
}

// initMessageSource 初始化消息源
func (ctx *DefaultApplicationContext) initMessageSource() error {
	// 消息源初始化逻辑
	return nil
}

// initApplicationEventMulticaster 初始化事件广播器
func (ctx *DefaultApplicationContext) initApplicationEventMulticaster() error {
	// 事件广播器初始化逻辑
	return nil
}

// registerListeners 注册监听器
func (ctx *DefaultApplicationContext) registerListeners() error {
	// 直接从工厂取，避免反rlock导致死锁
	listenersMap := ctx.beanFactory.GetBeansByType((*hdevcore.ApplicationListener)(nil))
	
	for _, listener := range listenersMap {
		if appListener, ok := listener.(hdevcore.ApplicationListener); ok {
			ctx.eventMulticaster.AddApplicationListener(appListener)
		}
	}
	
	// 发布早期事件
	ctx.publishEarlyEvents()
	
	return nil
}

// finishBeanFactoryInitialization 完成Bean工厂初始化
func (ctx *DefaultApplicationContext) finishBeanFactoryInitialization(beanFactory hdevcore.ConfigurableListableBeanFactory) error {
	// 预实例化所有单例Bean
	beanFactory.PreInstantiateSingletons()
	return nil
}

// finishRefresh 完成刷新
func (ctx *DefaultApplicationContext) finishRefresh() error {
	// 发布上下文刷新完成事件
	ctx.publishEvent(newContextRefreshedEvent(ctx))
	return nil
}

// publishEvent 发布事件
func (ctx *DefaultApplicationContext) publishEvent(event hdevcore.ApplicationEvent) {
	if ctx.eventMulticaster != nil {
		ctx.eventMulticaster.PublishEvent(event)
	}
}

// publishEarlyEvents 发布早期事件
func (ctx *DefaultApplicationContext) publishEarlyEvents() {
	for _, event := range ctx.earlyApplicationEvents {
		ctx.publishEvent(event)
	}
	ctx.earlyApplicationEvents = nil
}

// destroyBeans 销毁Bean
func (ctx *DefaultApplicationContext) destroyBeans() {
	// 销毁单例Bean的逻辑
}

// ContextRefreshedEvent 上下文刷新事件
type ContextRefreshedEvent struct {
	*hdevcore.GenericApplicationEvent
}

func newContextRefreshedEvent(source interface{}) *ContextRefreshedEvent {
	return &ContextRefreshedEvent{
		GenericApplicationEvent: hdevcore.NewGenericApplicationEvent(source),
	}
}

// ContextClosedEvent 上下文关闭事件
type ContextClosedEvent struct {
	*hdevcore.GenericApplicationEvent
}

func newContextClosedEvent(source interface{}) *ContextClosedEvent {
	return &ContextClosedEvent{
		GenericApplicationEvent: hdevcore.NewGenericApplicationEvent(source),
	}
}

// ContextStartedEvent 上下文启动事件
type ContextStartedEvent struct {
	*hdevcore.GenericApplicationEvent
}

// NewContextStartedEvent 创建上下文启动事件
func NewContextStartedEvent(source interface{}) *ContextStartedEvent {
	return &ContextStartedEvent{
		GenericApplicationEvent: hdevcore.NewGenericApplicationEvent(source),
	}
}

// ContextStoppedEvent 上下文停止事件
type ContextStoppedEvent struct {
	*hdevcore.GenericApplicationEvent
}

// NewContextStoppedEvent 创建上下文停止事件
func NewContextStoppedEvent(source interface{}) *ContextStoppedEvent {
	return &ContextStoppedEvent{
		GenericApplicationEvent: hdevcore.NewGenericApplicationEvent(source),
	}
}

// AddApplicationListener 添加应用事件监听器
func (ctx *DefaultApplicationContext) AddApplicationListener(listener hdevcore.ApplicationListener) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.eventMulticaster != nil {
		ctx.eventMulticaster.AddApplicationListener(listener)
	}
}

// RemoveApplicationListener 移除应用事件监听器
func (ctx *DefaultApplicationContext) RemoveApplicationListener(listener hdevcore.ApplicationListener) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.eventMulticaster != nil {
		ctx.eventMulticaster.RemoveApplicationListener(listener)
	}
}

// ApplicationContextBuilder 应用上下文构建器
type ApplicationContextBuilder struct {
	context *DefaultApplicationContext
}

// NewApplicationContextBuilder 创建应用上下文构建器
func NewApplicationContextBuilder() *ApplicationContextBuilder {
	return &ApplicationContextBuilder{
		context: NewDefaultApplicationContext(),
	}
}

// WithParent 设置父上下文
func (builder *ApplicationContextBuilder) WithParent(parent hdevcore.ApplicationContext) *ApplicationContextBuilder {
	builder.context.parent = parent
	return builder
}

// WithEnvironment 设置环境
func (builder *ApplicationContextBuilder) WithEnvironment(env hdevcore.Environment) *ApplicationContextBuilder {
	builder.context.SetEnvironment(env)
	return builder
}

// Build 构建应用上下文
func (builder *ApplicationContextBuilder) Build() *DefaultApplicationContext {
	return builder.context
}

// QuickStartApplicationContext 快速启动应用上下文
func QuickStartApplicationContext() (hdevcore.ApplicationContext, error) {
	ctx := NewDefaultApplicationContext()
	err := ctx.Refresh()
	return ctx, err
}