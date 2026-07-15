package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/zhifenghao123/how-dev-go-framework/hdev-config"
	"github.com/zhifenghao123/how-dev-go-framework/hdev-context"
	"github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// 测试用的Bean定义

type TestService struct {
	Config     *TestConfig
	Repository *TestRepository
	Initialized bool
	Destroyed   bool
}

type TestConfig struct {
	AppName string `yaml:"app.name"`
	Version string `yaml:"app.version"`
}

type TestRepository struct {
	DatabaseURL string
}

func (s *TestService) AfterPropertiesSet() error {
	s.Initialized = true
	return nil
}

func (s *TestService) Destroy() error {
	s.Destroyed = true
	return nil
}

func (s *TestService) GetInfo() string {
	return fmt.Sprintf("Service: %s v%s, DB: %s", s.Config.AppName, s.Config.Version, s.Repository.DatabaseURL)
}

// 测试事件监听器
type TestEventListener struct {
	Events []hdevcore.ApplicationEvent
}

func (l *TestEventListener) OnApplicationEvent(event hdevcore.ApplicationEvent) {
	l.Events = append(l.Events, event)
}

// 测试Bean后置处理器
type TestBeanPostProcessor struct {
	BeforeInitCount int
	AfterInitCount  int
}

func (p *TestBeanPostProcessor) PostProcessBeforeInitialization(bean interface{}, beanName string) (interface{}, error) {
	p.BeforeInitCount++
	return bean, nil
}

func (p *TestBeanPostProcessor) PostProcessAfterInitialization(bean interface{}, beanName string) (interface{}, error) {
	p.AfterInitCount++
	return bean, nil
}

// 集成测试函数

func TestFrameworkIntegration(t *testing.T) {
	t.Run("ApplicationContext Lifecycle", testApplicationContextLifecycle)
	t.Run("Configuration Binding", testConfigurationBinding)
	t.Run("Event System", testEventSystem)
	t.Run("Bean Factory", testBeanFactory)
	t.Run("Environment Integration", testEnvironmentIntegration)
	t.Run("Complete Workflow", testCompleteWorkflow)
}

func testApplicationContextLifecycle(t *testing.T) {
	// 创建应用上下文
	ctx := hdevcontext.NewDefaultApplicationContext()
	defer ctx.Close()

	// 注册Bean
	ctx.GetBeanFactory().RegisterSingleton("testConfig", &TestConfig{
		AppName: "Test App",
		Version: "1.0.0",
	})

	ctx.GetBeanFactory().RegisterSingleton("testRepository", &TestRepository{
		DatabaseURL: "test.db",
	})

	ctx.GetBeanFactory().RegisterSingleton("testService", &TestService{
		Config:     ctx.GetBeanFactory().GetBean("testConfig").(*TestConfig),
		Repository: ctx.GetBeanFactory().GetBean("testRepository").(*TestRepository),
	})

	// 刷新上下文
	if err := ctx.Refresh(); err != nil {
		t.Fatalf("Failed to refresh context: %v", err)
	}

	// 验证Bean初始化
	service := ctx.GetBean("testService").(*TestService)
	if !service.Initialized {
		t.Error("Service should be initialized")
	}

	// 验证依赖注入
	if service.Config.AppName != "Test App" {
		t.Errorf("Expected AppName 'Test App', got '%s'", service.Config.AppName)
	}

	if service.Repository.DatabaseURL != "test.db" {
		t.Errorf("Expected DatabaseURL 'test.db', got '%s'", service.Repository.DatabaseURL)
	}

	// 验证Bean功能
	expectedInfo := "Service: Test App v1.0.0, DB: test.db"
	if actualInfo := service.GetInfo(); actualInfo != expectedInfo {
		t.Errorf("Expected info '%s', got '%s'", expectedInfo, actualInfo)
	}
}

func testConfigurationBinding(t *testing.T) {
	// 创建配置环境
	env := hdevconfig.NewConfigurableEnvironment()

	// 添加配置源
	props := map[string]interface{}{
		"app.name":       "Config Test App",
		"app.version":    "2.0.0",
		"server.port":    "8080",
		"database.url":   "postgres://localhost:5432/test",
		"cache.enabled":  "true",
	}

	env.AddPropertySource(hdevconfig.NewMapPropertySource("test", props))

	// 创建配置绑定器
	binder := hdevconfig.NewConfigurationBinder(env)

	// 绑定配置
	var config TestConfig
	if err := binder.Bind("", &config); err != nil {
		t.Fatalf("Failed to bind configuration: %v", err)
	}

	// 验证配置绑定
	if config.AppName != "Config Test App" {
		t.Errorf("Expected AppName 'Config Test App', got '%s'", config.AppName)
	}

	if config.Version != "2.0.0" {
		t.Errorf("Expected Version '2.0.0', got '%s'", config.Version)
	}

	// 验证环境属性访问
	if appName := env.GetProperty("app.name"); appName != "Config Test App" {
		t.Errorf("Expected property 'Config Test App', got '%s'", appName)
	}

	if serverPort := env.GetProperty("server.port"); serverPort != "8080" {
		t.Errorf("Expected property '8080', got '%s'", serverPort)
	}

	// 验证属性存在性检查
	if !env.ContainsProperty("app.name") {
		t.Error("Property 'app.name' should exist")
	}

	if env.ContainsProperty("nonexistent.property") {
		t.Error("Property 'nonexistent.property' should not exist")
	}
}

func testEventSystem(t *testing.T) {
	// 创建应用上下文
	ctx := hdevcontext.NewDefaultApplicationContext()
	defer ctx.Close()

	// 创建事件监听器
	listener := &TestEventListener{}
	ctx.AddApplicationListener(listener)

	// 刷新上下文（会发布ContextRefreshedEvent）
	if err := ctx.Refresh(); err != nil {
		t.Fatalf("Failed to refresh context: %v", err)
	}

	// 验证事件发布
	if len(listener.Events) == 0 {
		t.Error("No events were published")
	}

	// 查找ContextRefreshedEvent
	var foundRefreshEvent bool
	for _, event := range listener.Events {
		if _, ok := event.(*hdevcontext.ContextRefreshedEvent); ok {
			foundRefreshEvent = true
			break
		}
	}

	if !foundRefreshEvent {
		t.Error("ContextRefreshedEvent should be published")
	}

	// 发布自定义事件
	customEvent := hdevcontext.NewContextStartedEvent(ctx)

	ctx.PublishEvent(customEvent)

	// 验证自定义事件
	var foundCustomEvent bool
	for _, event := range listener.Events {
		if _, ok := event.(*hdevcontext.ContextStartedEvent); ok {
			foundCustomEvent = true
			break
		}
	}

	if !foundCustomEvent {
		t.Error("ContextStartedEvent should be published")
	}
}

func testBeanFactory(t *testing.T) {
	// 创建Bean工厂
	factory := hdevcontext.NewDefaultListableBeanFactory()

	// 创建Bean定义
	configDef := hdevcontext.NewBeanDefinition("testConfig", reflect.TypeOf((*TestConfig)(nil)))
	configDef.SetScope("singleton")
	factory.RegisterBeanDefinition("testConfig", configDef)

	repoDef := hdevcontext.NewBeanDefinition("testRepository", reflect.TypeOf((*TestRepository)(nil)))
	repoDef.SetScope("singleton")
	factory.RegisterBeanDefinition("testRepository", repoDef)

	serviceDef := hdevcontext.NewBeanDefinition("testService", reflect.TypeOf((*TestService)(nil)))
	serviceDef.SetScope("singleton")
	serviceDef.SetInitMethodName("AfterPropertiesSet")
	serviceDef.SetDestroyMethodName("Destroy")
	factory.RegisterBeanDefinition("testService", serviceDef)

	// 添加Bean后置处理器
	postProcessor := &TestBeanPostProcessor{}
	factory.AddBeanPostProcessor(postProcessor)

	// 预实例化单例Bean
	if err := factory.PreInstantiateSingletons(); err != nil {
		t.Fatalf("Failed to pre-instantiate singletons: %v", err)
	}

	// 验证Bean实例化
	service := factory.GetBean("testService").(*TestService)
	if service == nil {
		t.Error("Service bean should not be nil")
	}

	// 验证Bean后置处理器调用
	if postProcessor.BeforeInitCount == 0 {
		t.Error("Bean post processor BeforeInit should be called")
	}

	if postProcessor.AfterInitCount == 0 {
		t.Error("Bean post processor AfterInit should be called")
	}

	// 验证Bean名称
	beanNames := factory.GetBeanDefinitionNames()
	expectedNames := []string{"testConfig", "testRepository", "testService"}
	if len(beanNames) != len(expectedNames) {
		t.Errorf("Expected %d bean names, got %d", len(expectedNames), len(beanNames))
	}

	// 销毁单例Bean
	factory.DestroySingletons()

	// 验证Bean销毁
	if !service.Destroyed {
		t.Error("Service should be destroyed")
	}
}

func testEnvironmentIntegration(t *testing.T) {
	// 创建配置环境
	env := hdevconfig.NewConfigurableEnvironment()

	// 添加多个配置源
	env.AddPropertySource(hdevconfig.NewMapPropertySource("source1", map[string]interface{}{
		"app.name": "Source1 App",
		"common.property": "value1",
	}))

	env.AddPropertySource(hdevconfig.NewMapPropertySource("source2", map[string]interface{}{
		"app.version": "3.0.0",
		"common.property": "value2", // 这个会覆盖source1的值
	}))

	// 验证配置源优先级（后添加的优先级高）
	if commonProp := env.GetProperty("common.property"); commonProp != "value2" {
		t.Errorf("Expected 'value2', got '%s'", commonProp)
	}

	// 验证属性合并
	if appName := env.GetProperty("app.name"); appName != "Source1 App" {
		t.Errorf("Expected 'Source1 App', got '%s'", appName)
	}

	if appVersion := env.GetProperty("app.version"); appVersion != "3.0.0" {
		t.Errorf("Expected '3.0.0', got '%s'", appVersion)
	}

	// 测试Profile功能
	env.SetActiveProfiles("development")

	env.AddPropertySource(hdevconfig.NewMapPropertySource("dev", map[string]interface{}{
		"profile.specific": "dev-value",
	}))

	if !env.AcceptsProfiles("development") {
		t.Error("Should accept development profile")
	}

	if profileSpecific := env.GetProperty("profile.specific"); profileSpecific != "dev-value" {
		t.Errorf("Expected 'dev-value', got '%s'", profileSpecific)
	}
}

func testCompleteWorkflow(t *testing.T) {
	// 完整的应用工作流测试
	
	// 1. 创建配置环境
	env := hdevconfig.NewConfigurableEnvironment()
	env.AddPropertySource(hdevconfig.NewMapPropertySource("app", map[string]interface{}{
		"app.name":       "Complete Test App",
		"app.version":    "1.0.0",
		"database.url":   "mysql://localhost:3306/test",
		"server.port":    "8080",
	}))

	// 2. 创建应用上下文
	ctx := hdevcontext.NewDefaultApplicationContext()
	ctx.SetEnvironment(env)
	defer ctx.Close()

	// 3. 添加事件监听器
	listener := &TestEventListener{}
	ctx.AddApplicationListener(listener)

	// 4. 添加Bean后置处理器
	postProcessor := &TestBeanPostProcessor{}
	ctx.GetBeanFactory().AddBeanPostProcessor(postProcessor)

	// 5. 注册Bean
	var config TestConfig
	binder := hdevconfig.NewConfigurationBinder(env)
	if err := binder.Bind("", &config); err != nil {
		t.Fatalf("Failed to bind configuration: %v", err)
	}

	ctx.GetBeanFactory().RegisterSingleton("appConfig", &config)
	ctx.GetBeanFactory().RegisterSingleton("testRepository", &TestRepository{
		DatabaseURL: env.GetProperty("database.url"),
	})
	ctx.GetBeanFactory().RegisterSingleton("testService", &TestService{
		Config:     ctx.GetBeanFactory().GetBean("appConfig").(*TestConfig),
		Repository: ctx.GetBeanFactory().GetBean("testRepository").(*TestRepository),
	})

	// 6. 刷新上下文
	if err := ctx.Refresh(); err != nil {
		t.Fatalf("Failed to refresh context: %v", err)
	}

	// 7. 验证完整功能
	service := ctx.GetBean("testService").(*TestService)
	
	// 验证Bean初始化
	if !service.Initialized {
		t.Error("Service should be initialized")
	}

	// 验证配置绑定
	if service.Config.AppName != "Complete Test App" {
		t.Errorf("Expected 'Complete Test App', got '%s'", service.Config.AppName)
	}

	// 验证依赖注入
	if service.Repository.DatabaseURL != "mysql://localhost:3306/test" {
		t.Errorf("Expected database URL, got '%s'", service.Repository.DatabaseURL)
	}

	// 验证事件发布
	if len(listener.Events) == 0 {
		t.Error("Events should be published")
	}

	// 验证Bean后置处理器
	if postProcessor.BeforeInitCount == 0 {
		t.Error("Bean post processor should be called")
	}

	// 发布业务事件
	ctx.PublishEvent(hdevcontext.NewContextStartedEvent(ctx))

	// 9. 关闭上下文
	ctx.Close()

	// 验证Bean销毁
	if !service.Destroyed {
		t.Error("Service should be destroyed")
	}
}

// 性能测试
func BenchmarkFrameworkInitialization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := hdevcontext.NewDefaultApplicationContext()
		
		// 注册一些Bean
		ctx.GetBeanFactory().RegisterSingleton("config", &TestConfig{AppName: "Benchmark"})
		ctx.GetBeanFactory().RegisterSingleton("service", &TestService{})
		
		ctx.Refresh()
		ctx.Close()
	}
}

// 主测试函数
func TestMain(m *testing.M) {
	fmt.Println("=== Starting How-Dev Framework Integration Tests ===")
	
	// 运行测试
	m.Run()
	
	fmt.Println("=== Integration Tests Completed ===")
}