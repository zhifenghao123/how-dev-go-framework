package hdev_ioc

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// BeanFactory Bean工厂接口
type BeanFactory interface {
	hdevcore.BeanFactory
	
	// RegisterBeanDefinition 注册Bean定义
	RegisterBeanDefinition(beanName string, beanDefinition *BeanDefinition) error
	// GetBeanDefinition 获取Bean定义
	GetBeanDefinition(beanName string) (*BeanDefinition, error)
	// ContainsBeanDefinition 检查是否包含Bean定义
	ContainsBeanDefinition(beanName string) bool
	// GetBeanDefinitionNames 获取所有Bean定义名称
	GetBeanDefinitionNames() []string
	// PreInstantiateSingletons 预实例化单例Bean
	PreInstantiateSingletons() error
}

// DefaultBeanFactory 默认Bean工厂实现
type DefaultBeanFactory struct {
	registry              BeanDefinitionRegistry
	singletonObjects      map[string]interface{}
	singletonMutex        sync.RWMutex
	beanPostProcessors   []hdevcore.BeanPostProcessor
	destroyed             bool
	context               context.Context
	cfg                   configProvider
	mu                    sync.RWMutex
}

// NewDefaultBeanFactory 创建默认Bean工厂
func NewDefaultBeanFactory() *DefaultBeanFactory {
	return &DefaultBeanFactory{
		registry:            NewDefaultBeanDefinitionRegistry(),
		singletonObjects:    make(map[string]interface{}),
		beanPostProcessors: make([]hdevcore.BeanPostProcessor, 0),
		destroyed:          false,
		context:            context.Background(),
	}
}

// Register 注册Bean实例
// 兼容旧式用法：当传入值类型 struct（如 Controller{}）时，自动转为指针实例 *Controller
// 这样后续依赖注入才能通过 reflect.Value.Set() 写入字段
func (factory *DefaultBeanFactory) Register(bean interface{}) {
	beanType := reflect.TypeOf(bean)
	if beanType == nil {
		return
	}
	
	// 值类型 struct → 转为对应指针类型的零值实例
	if beanType.Kind() == reflect.Struct {
		ptrInstance := reflect.New(beanType)
		// 复制原值到指针指向的位置（保留字段初始值）
		ptrInstance.Elem().Set(reflect.ValueOf(bean))
		bean = ptrInstance.Interface()
		beanType = reflect.TypeOf(bean)
	}
	
	beanName := factory.generateBeanName(beanType)
	
	beanDefinition := NewBeanDefinition(beanName, beanType)
	beanDefinition.SetInstance(bean)
	
	err := factory.registry.RegisterBeanDefinition(beanName, beanDefinition)
	if err != nil {
		panic(fmt.Sprintf("Failed to register bean: %v", err))
	}
}

// Instance 获取Bean实例
// 兼容多种调用方式：
//   - Instance("BeanName")：按名查找
//   - Instance(&MyType{}) 或 Instance(MyType{})：按类型查找
//   - Instance(reflect.Type)：按类型查找
func (factory *DefaultBeanFactory) Instance(bean interface{}) interface{} {
	if factory.destroyed {
		panic("BeanFactory has been destroyed")
	}
	
	// 字符串：按名称获取
	if name, ok := bean.(string); ok {
		return factory.GetBean(name)
	}
	
	// reflect.Type：直接使用
	var beanType reflect.Type
	if t, ok := bean.(reflect.Type); ok {
		beanType = t
	} else {
		beanType = reflect.TypeOf(bean)
	}
	if beanType == nil {
		return nil
	}
	
	beanName := factory.generateBeanName(beanType)
	
	// 首先尝试从单例缓存中获取
	if instance := factory.getSingleton(beanName); instance != nil {
		return instance
	}
	
	// 获取Bean定义
	beanDefinition, err := factory.registry.GetBeanDefinition(beanName)
	if err != nil {
		// 如果没有找到Bean定义，尝试创建新实例
		instance, err := factory.createBean(beanName, beanType)
		if err != nil {
			panic(fmt.Sprintf("Failed to create bean %s: %v", beanName, err))
		}
		return instance
	}
	
	// 根据作用域创建Bean实例
	instance, err := factory.getBean(beanName, beanDefinition)
	if err != nil {
		panic(fmt.Sprintf("Failed to get bean %s: %v", beanName, err))
	}
	return instance
}

// WithContext 设置上下文
func (factory *DefaultBeanFactory) WithContext(ctx context.Context) {
	factory.context = ctx
}

// RegisterBeanDefinition 注册Bean定义
func (factory *DefaultBeanFactory) RegisterBeanDefinition(beanName string, beanDefinition *BeanDefinition) error {
	return factory.registry.RegisterBeanDefinition(beanName, beanDefinition)
}

// GetBeanDefinition 获取Bean定义
func (factory *DefaultBeanFactory) GetBeanDefinition(beanName string) (*BeanDefinition, error) {
	return factory.registry.GetBeanDefinition(beanName)
}

// ContainsBeanDefinition 检查是否包含Bean定义
func (factory *DefaultBeanFactory) ContainsBeanDefinition(beanName string) bool {
	return factory.registry.ContainsBeanDefinition(beanName)
}

// GetBeanDefinitionNames 获取所有Bean定义名称
func (factory *DefaultBeanFactory) GetBeanDefinitionNames() []string {
	return factory.registry.GetBeanDefinitionNames()
}

// PreInstantiateSingletons 预实例化单例Bean
func (factory *DefaultBeanFactory) PreInstantiateSingletons() error {
	beanNames := factory.GetBeanDefinitionNames()
	
	for _, beanName := range beanNames {
		beanDefinition, err := factory.GetBeanDefinition(beanName)
		if err != nil {
			return err
		}
		
		// 如果是单例Bean且不是延迟初始化，则预实例化
		if beanDefinition.GetScope() == ScopeSingleton && !beanDefinition.IsLazyInit() {
			_, err := factory.getBean(beanName, beanDefinition)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

// AddBeanPostProcessor 添加Bean后置处理器
func (factory *DefaultBeanFactory) AddBeanPostProcessor(processor hdevcore.BeanPostProcessor) {
	factory.mu.Lock()
	defer factory.mu.Unlock()
	
	factory.beanPostProcessors = append(factory.beanPostProcessors, processor)
}

// GetBeanPostProcessors 获取Bean后置处理器
func (factory *DefaultBeanFactory) GetBeanPostProcessors() []hdevcore.BeanPostProcessor {
	factory.mu.RLock()
	defer factory.mu.RUnlock()
	
	return factory.beanPostProcessors
}

// Destroy 销毁Bean工厂
func (factory *DefaultBeanFactory) Destroy() {
	factory.mu.Lock()
	defer factory.mu.Unlock()
	
	if factory.destroyed {
		return
	}
	
	// 销毁单例Bean
	factory.singletonMutex.Lock()
	for beanName, bean := range factory.singletonObjects {
		if destroyable, ok := bean.(hdevcore.DisposableBean); ok {
			destroyable.Destroy()
		}
		delete(factory.singletonObjects, beanName)
	}
	factory.singletonMutex.Unlock()
	
	factory.destroyed = true
}

// getSingleton 获取单例Bean
func (factory *DefaultBeanFactory) getSingleton(beanName string) interface{} {
	factory.singletonMutex.RLock()
	defer factory.singletonMutex.RUnlock()
	
	return factory.singletonObjects[beanName]
}

// addSingleton 添加单例Bean
func (factory *DefaultBeanFactory) addSingleton(beanName string, singletonObject interface{}) {
	factory.singletonMutex.Lock()
	defer factory.singletonMutex.Unlock()
	
	factory.singletonObjects[beanName] = singletonObject
}

// getBean 获取Bean实例
func (factory *DefaultBeanFactory) getBean(beanName string, beanDefinition *BeanDefinition) (interface{}, error) {
	// 根据作用域处理
	switch beanDefinition.GetScope() {
	case ScopeSingleton:
		return factory.getSingletonBean(beanName, beanDefinition)
	case ScopePrototype:
		return factory.createBean(beanName, beanDefinition.GetBeanType())
	default:
		return nil, fmt.Errorf("unsupported bean scope: %s", beanDefinition.GetScope())
	}
}

// getSingletonBean 获取单例Bean
func (factory *DefaultBeanFactory) getSingletonBean(beanName string, beanDefinition *BeanDefinition) (interface{}, error) {
	// 双重检查锁定
	if singleton := factory.getSingleton(beanName); singleton != nil {
		return singleton, nil
	}
	
	// 1) 先在锁内确定/放入单例缓存（避免重复创建和循环依赖）
	var instance interface{}
	var newlyCreated bool
	
	factory.singletonMutex.Lock()
	if existing, ok := factory.singletonObjects[beanName]; ok {
		factory.singletonMutex.Unlock()
		return existing, nil
	}
	
	// 优先复用 BeanDefinition 中预设的实例（Register 注册的情况）
	if existing := beanDefinition.GetInstance(); existing != nil {
		instance = existing
	} else {
		// 创建空实例（不做字段注入），先放入单例缓存以打破循环依赖
		bt := beanDefinition.GetBeanType()
		if bt == nil {
			factory.singletonMutex.Unlock()
			return nil, fmt.Errorf("bean type is nil for %s", beanName)
		}
		if bt.Kind() == reflect.Ptr {
			instance = reflect.New(bt.Elem()).Interface()
		} else if bt.Kind() == reflect.Struct {
			instance = reflect.New(bt).Interface()
		} else {
			factory.singletonMutex.Unlock()
			return nil, fmt.Errorf("unsupported bean type kind: %s", bt.Kind())
		}
		newlyCreated = true
	}
	
	factory.singletonObjects[beanName] = instance
	factory.singletonMutex.Unlock()
	
	// 2) 在锁外做字段注入和初始化（允许内部递归取依赖Bean）
	defer func() {
		if r := recover(); r != nil {
			// 注入失败时清理缓存（避免空Bean被外部取走）
			factory.singletonMutex.Lock()
			delete(factory.singletonObjects, beanName)
			factory.singletonMutex.Unlock()
			panic(r)
		}
	}()
	
	// BeanPostProcessor: before
	if newlyCreated {
		instance = factory.applyBeanPostProcessorsBeforeInitialization(instance, beanName)
	}
	
	// 字段注入
	factory.injectFields(instance)
	
	// Init 钩子（兼容 BaseDao.Init(self)）
	factory.invokeInitHook(instance)
	
	// AfterPropertiesSet
	if initBean, ok := instance.(hdevcore.InitializingBean); ok {
		if err := initBean.AfterPropertiesSet(); err != nil {
			factory.singletonMutex.Lock()
			delete(factory.singletonObjects, beanName)
			factory.singletonMutex.Unlock()
			return nil, err
		}
	}
	
	// BeanPostProcessor: after
	if newlyCreated {
		instance = factory.applyBeanPostProcessorsAfterInitialization(instance, beanName)
		// 后置处理可能返回新对象，要回写缓存
		factory.singletonMutex.Lock()
		factory.singletonObjects[beanName] = instance
		factory.singletonMutex.Unlock()
	}
	
	return instance, nil
}

// createBean 创建Bean实例
func (factory *DefaultBeanFactory) createBean(beanName string, beanType reflect.Type) (interface{}, error) {
	// 创建Bean实例：要求 beanType 为指针类型；如果不是则自动取指针
	var beanValue interface{}
	if beanType.Kind() == reflect.Ptr {
		beanValue = reflect.New(beanType.Elem()).Interface()
	} else if beanType.Kind() == reflect.Struct {
		beanValue = reflect.New(beanType).Interface()
	} else {
		return nil, fmt.Errorf("unsupported bean type kind: %s", beanType.Kind())
	}
	
	// 应用Bean后置处理器
	beanValue = factory.applyBeanPostProcessorsBeforeInitialization(beanValue, beanName)
	
	// 字段注入（wired / value / const tag）
	factory.injectFields(beanValue)
	
	// 调用旧式的 Init(self) 钩子（兼容 hdev-gorm.BaseDao 等）
	factory.invokeInitHook(beanValue)
	
	// 初始化Bean
	if initializingBean, ok := beanValue.(hdevcore.InitializingBean); ok {
		if err := initializingBean.AfterPropertiesSet(); err != nil {
			return nil, err
		}
	}
	
	// 应用Bean后置处理器
	beanValue = factory.applyBeanPostProcessorsAfterInitialization(beanValue, beanName)
	
	return beanValue, nil
}

// invokeInitHook 调用 Init(self) 风格的初始化钩子
// 兼容 hdev-gorm.BaseDao.Init 等约定：如果 instance 拥有 Init 方法，且第一个参数与 instance 类型兼容，则调用之
func (factory *DefaultBeanFactory) invokeInitHook(instance interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// Init 钩子失败不影响Bean创建
		}
	}()
	v := reflect.ValueOf(instance)
	method := v.MethodByName("Init")
	if !method.IsValid() {
		return
	}
	mt := method.Type()
	switch mt.NumIn() {
	case 0:
		method.Call(nil)
	case 1:
		// 期望传入 self（IDao 风格）
		paramType := mt.In(0)
		if reflect.TypeOf(instance).AssignableTo(paramType) {
			method.Call([]reflect.Value{v})
		} else if paramType.Kind() == reflect.Interface && reflect.TypeOf(instance).Implements(paramType) {
			method.Call([]reflect.Value{v})
		}
	}
}

// applyBeanPostProcessorsBeforeInitialization 应用Bean后置处理器（初始化前）
func (factory *DefaultBeanFactory) applyBeanPostProcessorsBeforeInitialization(bean interface{}, beanName string) interface{} {
	result := bean
	
	for _, processor := range factory.GetBeanPostProcessors() {
		processedBean, err := processor.PostProcessBeforeInitialization(result, beanName)
		if err != nil || processedBean == nil {
			return result
		}
		result = processedBean
	}
	
	return result
}

// applyBeanPostProcessorsAfterInitialization 应用Bean后置处理器（初始化后）
func (factory *DefaultBeanFactory) applyBeanPostProcessorsAfterInitialization(bean interface{}, beanName string) interface{} {
	result := bean
	
	for _, processor := range factory.GetBeanPostProcessors() {
		processedBean, err := processor.PostProcessAfterInitialization(result, beanName)
		if err != nil || processedBean == nil {
			return result
		}
		result = processedBean
	}
	
	return result
}

// generateBeanName 生成Bean名称
func (factory *DefaultBeanFactory) generateBeanName(beanType reflect.Type) string {
	if beanType == nil {
		return ""
	}
	t := beanType
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// BeanFactoryBuilder Bean工厂构建器
type BeanFactoryBuilder struct {
	beanFactory *DefaultBeanFactory
}

// NewBeanFactoryBuilder 创建Bean工厂构建器
func NewBeanFactoryBuilder() *BeanFactoryBuilder {
	return &BeanFactoryBuilder{
		beanFactory: NewDefaultBeanFactory(),
	}
}

// WithBeanPostProcessor 添加Bean后置处理器
func (builder *BeanFactoryBuilder) WithBeanPostProcessor(processor hdevcore.BeanPostProcessor) *BeanFactoryBuilder {
	builder.beanFactory.AddBeanPostProcessor(processor)
	return builder
}

// WithContext 设置上下文
func (builder *BeanFactoryBuilder) WithContext(ctx context.Context) *BeanFactoryBuilder {
	builder.beanFactory.WithContext(ctx)
	return builder
}

// Build 构建Bean工厂
func (builder *BeanFactoryBuilder) Build() BeanFactory {
	return builder.beanFactory
}

// QuickStartBeanFactory 快速启动Bean工厂
func QuickStartBeanFactory() BeanFactory {
	return NewDefaultBeanFactory()
}

// GetBean 通过名称获取Bean
func (factory *DefaultBeanFactory) GetBean(name string) interface{} {
	if factory.destroyed {
		return nil
	}
	if instance := factory.getSingleton(name); instance != nil {
		return instance
	}
	beanDefinition, err := factory.registry.GetBeanDefinition(name)
	if err != nil {
		return nil
	}
	instance, err := factory.getBean(name, beanDefinition)
	if err != nil {
		return nil
	}
	return instance
}

// GetBeanByType 通过类型获取Bean
func (factory *DefaultBeanFactory) GetBeanByType(beanType interface{}) interface{} {
	return factory.Instance(beanType)
}

// ContainsBean 检查是否包含Bean
func (factory *DefaultBeanFactory) ContainsBean(name string) bool {
	if factory.getSingleton(name) != nil {
		return true
	}
	return factory.registry.ContainsBeanDefinition(name)
}

// IsSingleton 检查是否为单例
func (factory *DefaultBeanFactory) IsSingleton(name string) bool {
	beanDefinition, err := factory.registry.GetBeanDefinition(name)
	if err != nil {
		return false
	}
	return beanDefinition.GetScope() == ScopeSingleton
}

// IsPrototype 检查是否为原型
func (factory *DefaultBeanFactory) IsPrototype(name string) bool {
	beanDefinition, err := factory.registry.GetBeanDefinition(name)
	if err != nil {
		return false
	}
	return beanDefinition.GetScope() == ScopePrototype
}

// GetType 获取Bean类型
func (factory *DefaultBeanFactory) GetType(name string) interface{} {
	beanDefinition, err := factory.registry.GetBeanDefinition(name)
	if err != nil {
		return nil
	}
	return beanDefinition.GetBeanType()
}

// GetAliases 获取别名
func (factory *DefaultBeanFactory) GetAliases(name string) []string {
	return []string{}
}

// Example usage:
/*
type UserService struct {
	UserRepository *UserRepository
}

func (s *UserService) AfterPropertiesSet() error {
	// 初始化逻辑
	return nil
}

func main() {
	// 创建Bean工厂
	beanFactory := hdevioc.NewDefaultBeanFactory()
	
	// 注册Bean
	beanFactory.Register(&UserService{})
	
	// 预实例化单例Bean
	beanFactory.PreInstantiateSingletons()
	
	// 获取Bean实例
	userService := beanFactory.Instance(&UserService{}).(*UserService)
	
	fmt.Printf("UserService: %+v\n", userService)
}
*/