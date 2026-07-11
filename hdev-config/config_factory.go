package hdevconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// ConfigFactory 配置工厂
type ConfigFactory struct {
	baseDir        string
	activeProfiles []string
	configFiles    []string
}

// NewConfigFactory 创建配置工厂
func NewConfigFactory() *ConfigFactory {
	return &ConfigFactory{
		baseDir:        ".",
		activeProfiles: make([]string, 0),
		configFiles:    make([]string, 0),
	}
}

// WithBaseDir 设置基础目录
func (factory *ConfigFactory) WithBaseDir(dir string) *ConfigFactory {
	factory.baseDir = dir
	return factory
}

// WithActiveProfiles 设置激活的环境
func (factory *ConfigFactory) WithActiveProfiles(profiles ...string) *ConfigFactory {
	factory.activeProfiles = profiles
	return factory
}

// WithConfigFiles 设置配置文件
func (factory *ConfigFactory) WithConfigFiles(files ...string) *ConfigFactory {
	factory.configFiles = files
	return factory
}

// CreateEnvironment 创建配置环境
func (factory *ConfigFactory) CreateEnvironment() (hdevcore.Environment, error) {
	env := NewConfigurableEnvironment()
	
	// 设置激活的环境
	if len(factory.activeProfiles) > 0 {
		env.SetActiveProfiles(factory.activeProfiles...)
	}
	
	// 加载配置文件
	err := factory.loadConfigFiles(env)
	if err != nil {
		return nil, err
	}
	
	return env, nil
}

// loadConfigFiles 加载配置文件
func (factory *ConfigFactory) loadConfigFiles(env hdevcore.Environment) error {
	// 如果没有指定配置文件，使用默认的配置文件查找策略
	if len(factory.configFiles) == 0 {
		return factory.loadDefaultConfigFiles(env)
	}
	
	// 加载指定的配置文件
	for _, file := range factory.configFiles {
		err := factory.loadConfigFile(env, file)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// loadDefaultConfigFiles 加载默认配置文件
func (factory *ConfigFactory) loadDefaultConfigFiles(env hdevcore.Environment) error {
	// 默认配置文件查找策略：
	// 1. application.{ext}
	// 2. application-{profile}.{ext}
	// 3. 环境特定的配置文件
	
	extensions := []ConfigFormat{ConfigFormatYAML, ConfigFormatJSON, ConfigFormatProperties}
	
	// 加载基础配置文件
	for _, ext := range extensions {
		fileName := fmt.Sprintf("application.%s", ext)
		filePath := filepath.Join(factory.baseDir, fileName)
		
		if factory.fileExists(filePath) {
			err := factory.loadConfigFile(env, filePath)
			if err != nil {
				return err
			}
		}
	}
	
	// 加载环境特定的配置文件
	for _, profile := range factory.activeProfiles {
		for _, ext := range extensions {
			fileName := fmt.Sprintf("application-%s.%s", profile, ext)
			filePath := filepath.Join(factory.baseDir, fileName)
			
			if factory.fileExists(filePath) {
				err := factory.loadConfigFile(env, filePath)
				if err != nil {
					return err
				}
			}
		}
	}
	
	return nil
}

// loadConfigFile 加载单个配置文件
func (factory *ConfigFactory) loadConfigFile(env hdevcore.Environment, filePath string) error {
	// 确定文件格式
	ext := factory.getFileExtension(filePath)
	format := factory.getConfigFormat(ext)
	
	if format == "" {
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
	
	// 创建文件属性源
	sourceName := fmt.Sprintf("fileProperties[%s]", filepath.Base(filePath))
	source, err := NewFilePropertySource(sourceName, filePath, format)
	if err != nil {
		return err
	}
	
	// 添加到环境
	if configurableEnv, ok := env.(*ConfigurableEnvironment); ok {
		configurableEnv.AddPropertySource(source)
	}
	
	return nil
}

// getFileExtension 获取文件扩展名
func (factory *ConfigFactory) getFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}
	return ""
}

// getConfigFormat 根据扩展名获取配置格式
func (factory *ConfigFactory) getConfigFormat(ext string) ConfigFormat {
	switch strings.ToLower(ext) {
	case "yml", "yaml":
		return ConfigFormatYAML
	case "json":
		return ConfigFormatJSON
	case "ini":
		return ConfigFormatINI
	case "properties", "props":
		return ConfigFormatProperties
	default:
		return ""
	}
}

// fileExists 检查文件是否存在
func (factory *ConfigFactory) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DefaultConfigFactory 默认配置工厂
func DefaultConfigFactory() *ConfigFactory {
	factory := NewConfigFactory()
	
	// 从环境变量获取激活的环境
	if profiles := os.Getenv("SPRING_PROFILES_ACTIVE"); profiles != "" {
		factory.WithActiveProfiles(strings.Split(profiles, ",")...)
	}
	
	return factory
}

// QuickStart 快速启动方法
func QuickStart() (hdevcore.Environment, error) {
	return DefaultConfigFactory().CreateEnvironment()
}

// QuickStartWithProfiles 带环境配置的快速启动
func QuickStartWithProfiles(profiles ...string) (hdevcore.Environment, error) {
	return NewConfigFactory().WithActiveProfiles(profiles...).CreateEnvironment()
}

// QuickStartWithFiles 带配置文件的快速启动
func QuickStartWithFiles(files ...string) (hdevcore.Environment, error) {
	return NewConfigFactory().WithConfigFiles(files...).CreateEnvironment()
}

// Example usage:
/*
func main() {
	// 方式1：快速启动（使用默认配置）
	env, err := hdevconfig.QuickStart()
	if err != nil {
		log.Fatal(err)
	}
	
	// 方式2：指定环境
	env, err = hdevconfig.QuickStartWithProfiles("dev", "local")
	
	// 方式3：指定配置文件
	env, err = hdevconfig.QuickStartWithFiles("config/app.yaml", "config/db.properties")
	
	// 方式4：使用工厂模式
	factory := hdevconfig.NewConfigFactory().
		WithBaseDir("./config").
		WithActiveProfiles("prod").
		WithConfigFiles("app.yaml", "database.yaml")
	
	env, err = factory.CreateEnvironment()
	
	// 使用环境
	serverPort := env.GetPropertyWithDefault("server.port", "8080")
	fmt.Printf("Server will run on port %s\n", serverPort)
}
*/