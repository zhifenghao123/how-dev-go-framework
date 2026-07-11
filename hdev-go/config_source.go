package hdevgo

import (
	"fmt"
	"strings"

	hdevconfig "github.com/zhifenghao123/how-dev-go-framework/hdev-config"
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// ConfigPropertySource 把 hdev-config.Config 包装为 hdev-core.PropertySource
// 同时也实现了 configProvider（拥有 Get 方法），可被 hdev-ioc.DefaultBeanFactory.SetConfig 接受
type ConfigPropertySource struct {
	name string
	cfg  *hdevconfig.Config
}

// NewConfigPropertySource 创建配置属性源
func NewConfigPropertySource(name string, cfg *hdevconfig.Config) *ConfigPropertySource {
	return &ConfigPropertySource{name: name, cfg: cfg}
}

// GetName 获取属性源名称
func (s *ConfigPropertySource) GetName() string {
	return s.name
}

// GetProperty 获取属性
func (s *ConfigPropertySource) GetProperty(key string) interface{} {
	if s.cfg == nil {
		return nil
	}
	return s.cfg.Get(key)
}

// ContainsProperty 检查属性是否存在
func (s *ConfigPropertySource) ContainsProperty(key string) bool {
	if s.cfg == nil {
		return false
	}
	return s.cfg.Get(key) != nil
}

// Get 透传 hdev-config.IConfig 的 Get 方法（适配 hdev-ioc 的 configProvider 接口）
func (s *ConfigPropertySource) Get(key string, callbacks ...func(interface{})) interface{} {
	if s.cfg == nil {
		return nil
	}
	return s.cfg.Get(key, callbacks...)
}

// Config 返回底层 *hdevconfig.Config，供需要的组件直接使用
func (s *ConfigPropertySource) Config() *hdevconfig.Config {
	return s.cfg
}

// LoadConfigFile 加载主配置文件并自动处理 import 字段（导入其它 yaml）
// fileName 为不带后缀的文件名（如 "application"），返回内置的 ConfigPropertySource
func LoadConfigFile(file string) (*ConfigPropertySource, error) {
	cfg, err := hdevconfig.Init(hdevconfig.Option{File: file})
	if err != nil {
		return nil, fmt.Errorf("load config file %s error: %w", file, err)
	}

	cname := configName(file)

	// 处理 import 字段：${cname}.import: [path, path@alias]
	if list := cfg.GetAsArray(joinDot(cname, configKeyImport)); list != nil {
		for _, item := range list {
			s, _ := item.(string)
			if s == "" {
				continue
			}
			path, alias := extractImport(s)
			if err := cfg.Load(path, alias...); err != nil {
				return nil, fmt.Errorf("import config %s error: %w", s, err)
			}
		}
	}

	source := NewConfigPropertySource(fmt.Sprintf("config[%s]", file), cfg)
	return source, nil
}

// configName 提取去掉路径与扩展名后的文件名（与 hdev-config.Filename 保持一致语义）
func configName(file string) string {
	name := file
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}
	return name
}

// joinDot 用 . 连接 key 片段
func joinDot(parts ...string) string {
	return strings.Join(parts, ".")
}

// extractImport 解析形如 "path/to/file.yml@alias" 的引用
func extractImport(item string) (string, []string) {
	if !strings.Contains(item, "@") {
		return item, nil
	}
	arr := strings.Split(item, "@")
	return strings.TrimSpace(arr[0]), []string{strings.TrimSpace(arr[1])}
}

// 编译期接口断言
var _ hdevcore.PropertySource = (*ConfigPropertySource)(nil)
