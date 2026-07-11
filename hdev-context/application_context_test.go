package hdevcontext

import (
	"reflect"
	"testing"
	"time"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// TestUserService 测试用户服务
type TestUserService struct {
	Name string `autowired:""`
}

// TestUserRepository 测试用户仓库
type TestUserRepository struct {
	Database string
}

// TestConfiguration 测试配置类
type TestConfiguration struct{}

// TestBean 测试Bean方法
func (config *TestConfiguration) TestBean() *TestUserRepository {
	return &TestUserRepository{Database: "test_db"}
}

// TestApplicationListener 测试应用监听器
type TestApplicationListener struct {
	eventCount int
}

// OnApplicationEvent 处理应用事件
func (listener *TestApplicationListener) OnApplicationEvent(event hdevcore.ApplicationEvent) {
	listener.eventCount++
}

// TestApplicationContext 测试应用上下文
func TestApplicationContext(t *testing.T) {
	// 创建应用上下文
	ctx := NewDefaultApplicationContext()
	
	// 注册Bean定义
	userServiceDef := NewBeanDefinition("userService", reflect.TypeOf((*TestUserService)(nil)))
	userRepoDef := NewBeanDefinition("userRepository", reflect.TypeOf((*TestUserRepository)(nil)))
	
	ctx.GetBeanFactory().RegisterBeanDefinition("userService", userServiceDef)
	ctx.GetBeanFactory().RegisterBeanDefinition("userRepository", userRepoDef)
	
	// 刷新上下文
	err := ctx.Refresh()
	if err != nil {
		t.Fatalf("Failed to refresh context: %v", err)
	}
	
	// 测试获取Bean
	userService := ctx.GetBean("userService").(*TestUserService)
	if userService == nil {
		t.Error("Failed to get userService bean")
	}
	
	// 测试环境配置
	env := ctx.GetEnvironment()
	if env == nil {
		t.Error("Failed to get environment")
	}
	
	// 测试上下文状态
	if !ctx.IsActive() {
		t.Error("Context should be active after refresh")
	}
	
	// 关闭上下文
	err = ctx.Close()
	if err != nil {
		t.Errorf("Failed to close context: %v", err)
	}
	
	if ctx.IsActive() {
		t.Error("Context should not be active after close")
	}
}

// TestBeanFactory 测试Bean工厂
func TestBeanFactory(t *testing.T) {
	factory := NewDefaultListableBeanFactory()
	
	// 注册单例Bean
	factory.RegisterSingleton("testBean", &TestUserService{Name: "test"})
	
	// 测试Bean存在性
	if !factory.ContainsBean("testBean") {
		t.Error("Bean should exist after registration")
	}
	
	// 测试获取Bean
	bean := factory.GetBean("testBean")
	if bean == nil {
		t.Error("Failed to get bean")
	}
	
	// 测试单例特性
	bean2 := factory.GetBean("testBean")
	if bean != bean2 {
		t.Error("Singleton beans should be the same instance")
	}
	
	// 测试类型获取
	beanType := factory.GetType("testBean")
	if beanType == nil {
		t.Error("Failed to get bean type")
	}
}

// TestEnvironment 测试环境配置
func TestEnvironment(t *testing.T) {
	env := NewStandardEnvironment()
	
	// 测试属性源
	source := NewMapPropertySource("test", map[string]string{
		"app.name": "test-app",
		"app.version": "1.0.0",
	})
	env.AddPropertySource(source)
	
	// 测试属性获取
	name := env.GetProperty("app.name")
	if name != "test-app" {
		t.Errorf("Expected 'test-app', got '%s'", name)
	}
	
	// 测试默认值
	port := env.GetPropertyWithDefault("server.port", "8080")
	if port != "8080" {
		t.Errorf("Expected '8080', got '%s'", port)
	}
	
	// 测试属性存在性
	if !env.ContainsProperty("app.name") {
		t.Error("Property should exist")
	}
	
	// 测试profile
	env.SetActiveProfiles("test", "dev")
	if !env.AcceptsProfiles("test") {
		t.Error("Should accept test profile")
	}
}

// TestEventSystem 测试事件系统
func TestEventSystem(t *testing.T) {
	multicaster := NewSimpleApplicationEventMulticaster()
	listener := &TestApplicationListener{}
	
	// 注册监听器
	multicaster.AddApplicationListener(listener)
	
	// 发布事件
	event := hdevcore.NewGenericApplicationEvent("test")
	multicaster.PublishEvent(event)
	
	// 等待事件处理
	time.Sleep(100 * time.Millisecond)
	
	if listener.eventCount == 0 {
		t.Error("Listener should have received event")
	}
	
	// 测试监听器移除
	multicaster.RemoveApplicationListener(listener)
	oldCount := listener.eventCount
	
	multicaster.PublishEvent(event)
	time.Sleep(100 * time.Millisecond)
	
	if listener.eventCount != oldCount {
		t.Error("Listener should not receive events after removal")
	}
}

// TestBeanLifecycle 测试Bean生命周期
func TestBeanLifecycle(t *testing.T) {
	factory := NewDefaultListableBeanFactory()
	
	// 创建Bean定义
	def := NewBeanDefinition("lifecycleBean", reflect.TypeOf((*TestUserService)(nil)))
	def.SetInitMethodName("Initialize")
	
	factory.RegisterBeanDefinition("lifecycleBean", def)
	
	// 预实例化单例
	err := factory.PreInstantiateSingletons()
	if err != nil {
		t.Fatalf("Failed to pre-instantiate singletons: %v", err)
	}
	
	// 测试Bean获取
	bean := factory.GetBean("lifecycleBean")
	if bean == nil {
		t.Error("Failed to get lifecycle bean")
	}
}

// TestDependencyInjection 测试依赖注入
func TestDependencyInjection(t *testing.T) {
	t.Skip("依赖注入测试需要重新设计：TestUserService中的Name字段是string类型，无法通过类型注入TestUserRepository")
}

// TestApplicationContextBuilder 测试应用上下文构建器
func TestApplicationContextBuilder(t *testing.T) {
	builder := NewApplicationContextBuilder()
	
	// 构建上下文
	ctx := builder.Build()
	
	// 测试刷新
	err := ctx.Refresh()
	if err != nil {
		t.Fatalf("Failed to refresh context: %v", err)
	}
	
	// 测试基本功能
	if !ctx.IsActive() {
		t.Error("Context should be active")
	}
	
	// 测试关闭
	err = ctx.Close()
	if err != nil {
		t.Errorf("Failed to close context: %v", err)
	}
}

// BenchmarkBeanCreation 性能测试：Bean创建
func BenchmarkBeanCreation(b *testing.B) {
	factory := NewDefaultListableBeanFactory()
	
	// 注册测试Bean
	def := NewBeanDefinition("benchmarkBean", reflect.TypeOf((*TestUserService)(nil)))
	factory.RegisterBeanDefinition("benchmarkBean", def)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		factory.GetBean("benchmarkBean")
	}
}

// BenchmarkContextRefresh 性能测试：上下文刷新
func BenchmarkContextRefresh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := NewDefaultApplicationContext()
		ctx.Refresh()
		ctx.Close()
	}
}