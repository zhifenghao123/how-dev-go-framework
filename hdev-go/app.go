package hdevgo

import (
	"context"
	"log"
	"os"
	"reflect"
	"sync"

	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
	hdevioc "github.com/zhifenghao123/how-dev-go-framework/hdev-ioc"
)

// AppOption 应用启动选项
type AppOption func(*App)

// WithConf 指定主配置文件路径
func WithConf(file string) AppOption {
	return func(a *App) {
		a.confFile = file
	}
}

// WithDebug 开启调试模式（注入到 ctx 中，BaseDao 等会读取）
func WithDebug(debug bool) AppOption {
	return func(a *App) {
		a.debug = debug
	}
}

// App 应用启动器
// 提供 Spring Boot 风格的链式 API：NewApp(opts).Process(...).Register(...).Run()
type App struct {
	// 配置
	confFile string
	debug    bool

	// 运行时上下文
	ctx    context.Context
	cancel context.CancelFunc

	// 配置源（同时实现 PropertySource 与 hdev-ioc 的 configProvider）
	cfg *ConfigPropertySource
	// 主配置名（不带扩展名），用于拼接 process 配置 key
	cname string

	// IoC 容器（复用 hdev-ioc 的 BeanFactory，符合方案 1C）
	beanFactory *hdevioc.DefaultBeanFactory

	// 进程的 Bean 名称集合（按原型作用域注册）
	processBeans []string

	// 等待启动完成
	wg sync.WaitGroup
}

// NewApp 创建应用启动器
func NewApp(opts ...AppOption) *App {
	a := &App{
		confFile: "conf/application.yml",
	}
	for _, opt := range opts {
		opt(a)
	}

	// 初始化运行时 context
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.ctx = context.WithValue(a.ctx, "debug", a.debug)

	// 加载主配置（包含 import 处理）
	cfg, err := LoadConfigFile(a.confFile)
	if err != nil {
		log.Fatalf("[hdev-context] load config file error: %v", err)
	}
	a.cfg = cfg
	a.cname = configName(a.confFile)

	// 触发插件加载配置
	loadPlugins(a.cfg)

	// 初始化 BeanFactory
	a.beanFactory = hdevioc.NewDefaultBeanFactory()
	a.beanFactory.WithContext(a.ctx)
	a.beanFactory.SetConfig(a.cfg)

	// 登记为全局 App，供 hdev-gorm 等通过 hdevgo.Config(...) 访问
	setGlobalApp(a)

	return a
}

// Process 注册一个或多个 Process Bean（按原型作用域）
// Process Bean 通常实现 Execute() 方法，由 Run() 在启动时调度
func (a *App) Process(beans ...interface{}) *App {
	for _, b := range beans {
		// 先按实例注册（构建 BeanDefinition）
		a.beanFactory.Register(b)

		// 把作用域改为原型，避免与同名单例 Bean 冲突
		beanType := reflect.TypeOf(b)
		var beanName string
		if beanType.Kind() == reflect.Ptr {
			beanName = beanType.Elem().Name()
		} else {
			beanName = beanType.Name()
		}
		if bd, err := a.beanFactory.GetBeanDefinition(beanName); err == nil && bd != nil {
			bd.SetScope(hdevioc.ScopePrototype)
		}
		a.processBeans = append(a.processBeans, beanName)
	}
	return a
}

// Register 注册业务 Bean（Controller/Service/Dao/Middleware 等，单例）
// 支持以下入参形式：
//   - struct 实例（值或指针）
//   - []interface{} / []any 切片（自动展开）
func (a *App) Register(beans ...interface{}) *App {
	for _, b := range beans {
		if b == nil {
			continue
		}
		switch arr := b.(type) {
		case []interface{}:
			for _, v := range arr {
				a.beanFactory.Register(v)
			}
		default:
			a.beanFactory.Register(b)
		}
	}
	return a
}

// BeanFactory 暴露内部 BeanFactory，供高级用例（自定义初始化、二次注册等）使用
func (a *App) BeanFactory() *hdevioc.DefaultBeanFactory {
	return a.beanFactory
}

// Context 暴露运行时 context（监听取消信号）
func (a *App) Context() context.Context {
	return a.ctx
}

// Config 暴露配置源
func (a *App) Config() *ConfigPropertySource {
	return a.cfg
}

// Run 启动应用：信号监听 → 插件 Setup → 调度 Process → 等待结束 → 关闭插件
func (a *App) Run() {
	logPid()

	// 监听停止信号
	stopChan := make(chan bool, 1)
	go listenStopSignal(stopChan)
	go func() {
		<-stopChan
		log.Println("[hdev-context] stopping...")
		stopPlugins()
		a.cancel()
	}()

	// 启动前：插件 Setup
	setupPlugins()

	// 提取 Process 配置并启动
	processes := extractProcesses(a.cfg, a.cname)
	if len(processes) == 0 {
		log.Println("[hdev-context] no process found in config, exit")
	}
	for _, p := range processes {
		a.wg.Add(1)
		go runProcess(p, func(name string) interface{} {
			return a.beanFactory.GetBean(name)
		}, &a.wg)
	}
	a.wg.Wait()

	// 关闭插件
	if closePlugins() {
		_ = removePid()
		log.Println("[hdev-context] stop success. Goodbye.")
	}
}

// Close 主动关闭应用（不等待 Process 结束）
func (a *App) Close() {
	a.cancel()
	stopPlugins()
	closePlugins()
	_ = removePid()
}

// 编译期断言：BeanFactory 实现 hdev-core.BeanFactory 接口
var _ hdevcore.BeanFactory = (*hdevioc.DefaultBeanFactory)(nil)

// 防止 os 包导入未使用
var _ = os.Getpid
