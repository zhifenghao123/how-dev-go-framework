package hdevconfig

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// FilePropertySource 文件属性源
type FilePropertySource struct {
	name     string
	filePath string
	format   ConfigFormat
	data     map[string]interface{}
	mu       sync.RWMutex
}

// ConfigFormat 配置文件格式
type ConfigFormat string

const (
	ConfigFormatYAML ConfigFormat = "yaml"
	ConfigFormatJSON ConfigFormat = "json"
	ConfigFormatINI  ConfigFormat = "ini"
	ConfigFormatProperties ConfigFormat = "properties"
)

// NewFilePropertySource 创建文件属性源
func NewFilePropertySource(name, filePath string, format ConfigFormat) (*FilePropertySource, error) {
	source := &FilePropertySource{
		name:     name,
		filePath: filePath,
		format:   format,
		data:     make(map[string]interface{}),
	}
	
	err := source.load()
	if err != nil {
		return nil, err
	}
	
	return source, nil
}

// GetName 获取属性源名称
func (fps *FilePropertySource) GetName() string {
	return fps.name
}

// GetProperty 获取属性值
func (fps *FilePropertySource) GetProperty(key string) interface{} {
	fps.mu.RLock()
	defer fps.mu.RUnlock()
	
	// 支持嵌套属性访问，如 "server.port"
	keys := strings.Split(key, ".")
	current := fps.data
	
	for i, k := range keys {
		if i == len(keys)-1 {
			return current[k]
		}
		
		if next, ok := current[k].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}
	
	return nil
}

// ContainsProperty 检查是否包含指定属性
func (fps *FilePropertySource) ContainsProperty(key string) bool {
	fps.mu.RLock()
	defer fps.mu.RUnlock()
	
	keys := strings.Split(key, ".")
	current := fps.data
	
	for i, k := range keys {
		if i == len(keys)-1 {
			_, exists := current[k]
			return exists
		}
		
		if next, ok := current[k].(map[string]interface{}); ok {
			current = next
		} else {
			return false
		}
	}
	
	return false
}

// SetProperty 设置属性值
func (fps *FilePropertySource) SetProperty(key string, value interface{}) {
	fps.mu.Lock()
	defer fps.mu.Unlock()
	
	keys := strings.Split(key, ".")
	current := fps.data
	
	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
			return
		}
		
		if next, ok := current[k].(map[string]interface{}); ok {
			current = next
		} else {
			// 创建新的嵌套map
			nextMap := make(map[string]interface{})
			current[k] = nextMap
			current = nextMap
		}
	}
}

// Reload 重新加载配置文件
func (fps *FilePropertySource) Reload() error {
	fps.mu.Lock()
	defer fps.mu.Unlock()
	
	return fps.load()
}

// load 加载配置文件
func (fps *FilePropertySource) load() error {
	// 读取文件内容
	content, err := ioutil.ReadFile(fps.filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", fps.filePath, err)
	}
	
	// 根据格式解析文件内容
	parser, err := fps.getParser()
	if err != nil {
		return err
	}
	
	data, err := parser.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse config file %s: %v", fps.filePath, err)
	}
	
	fps.data = data
	return nil
}

// getParser 获取对应的解析器
func (fps *FilePropertySource) getParser() (ConfigParser, error) {
	switch fps.format {
	case ConfigFormatYAML:
		return &YAMLConfigParser{}, nil
	case ConfigFormatJSON:
		return &JSONConfigParser{}, nil
	case ConfigFormatINI:
		return &INIConfigParser{}, nil
	case ConfigFormatProperties:
		return &PropertiesConfigParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported config format: %s", fps.format)
	}
}

// ConfigParser 配置解析器接口
type ConfigParser interface {
	Parse(content []byte) (map[string]interface{}, error)
}

// YAMLConfigParser YAML配置解析器
type YAMLConfigParser struct{}

func (p *YAMLConfigParser) Parse(content []byte) (map[string]interface{}, error) {
	// 简化实现，实际应该使用yaml库解析
	return make(map[string]interface{}), nil
}

// JSONConfigParser JSON配置解析器
type JSONConfigParser struct{}

func (p *JSONConfigParser) Parse(content []byte) (map[string]interface{}, error) {
	// 简化实现，实际应该使用json库解析
	return make(map[string]interface{}), nil
}

// INIConfigParser INI配置解析器
type INIConfigParser struct{}

func (p *INIConfigParser) Parse(content []byte) (map[string]interface{}, error) {
	// 简化实现，实际应该使用ini库解析
	return make(map[string]interface{}), nil
}

// PropertiesConfigParser Properties配置解析器
type PropertiesConfigParser struct{}

func (p *PropertiesConfigParser) Parse(content []byte) (map[string]interface{}, error) {
	// 简化实现，实际应该解析properties格式
	return make(map[string]interface{}), nil
}

// CompositeFilePropertySource 组合文件属性源
type CompositeFilePropertySource struct {
	name            string
	propertySources []hdevcore.PropertySource
	mu              sync.RWMutex
}

// NewCompositeFilePropertySource 创建组合文件属性源
func NewCompositeFilePropertySource(name string) *CompositeFilePropertySource {
	return &CompositeFilePropertySource{
		name:            name,
		propertySources: make([]hdevcore.PropertySource, 0),
	}
}

func (cfps *CompositeFilePropertySource) GetName() string {
	return cfps.name
}

func (cfps *CompositeFilePropertySource) GetProperty(key string) interface{} {
	cfps.mu.RLock()
	defer cfps.mu.RUnlock()
	
	for i := len(cfps.propertySources) - 1; i >= 0; i-- {
		if value := cfps.propertySources[i].GetProperty(key); value != nil {
			return value
		}
	}
	return nil
}

func (cfps *CompositeFilePropertySource) ContainsProperty(key string) bool {
	cfps.mu.RLock()
	defer cfps.mu.RUnlock()
	
	for _, ps := range cfps.propertySources {
		if ps.ContainsProperty(key) {
			return true
		}
	}
	return false
}

// AddPropertySource 添加属性源
func (cfps *CompositeFilePropertySource) AddPropertySource(propertySource hdevcore.PropertySource) {
	cfps.mu.Lock()
	defer cfps.mu.Unlock()
	
	cfps.propertySources = append(cfps.propertySources, propertySource)
}

// ProfileSpecificFilePropertySource 特定环境文件属性源
type ProfileSpecificFilePropertySource struct {
	baseSource     hdevcore.PropertySource
	profileSources map[string]hdevcore.PropertySource
	mu             sync.RWMutex
}

// NewProfileSpecificFilePropertySource 创建特定环境文件属性源
func NewProfileSpecificFilePropertySource(baseSource hdevcore.PropertySource) *ProfileSpecificFilePropertySource {
	return &ProfileSpecificFilePropertySource{
		baseSource:     baseSource,
		profileSources: make(map[string]hdevcore.PropertySource),
	}
}

func (psfps *ProfileSpecificFilePropertySource) GetName() string {
	return psfps.baseSource.GetName()
}

func (psfps *ProfileSpecificFilePropertySource) GetProperty(key string) interface{} {
	psfps.mu.RLock()
	defer psfps.mu.RUnlock()
	
	// 先检查基础属性源
	if value := psfps.baseSource.GetProperty(key); value != nil {
		return value
	}
	
	// 这里应该根据当前激活的环境检查特定环境的属性源
	// 简化实现
	return nil
}

func (psfps *ProfileSpecificFilePropertySource) ContainsProperty(key string) bool {
	psfps.mu.RLock()
	defer psfps.mu.RUnlock()
	
	if psfps.baseSource.ContainsProperty(key) {
		return true
	}
	
	// 简化实现
	return false
}

// AddProfileSource 添加特定环境属性源
func (psfps *ProfileSpecificFilePropertySource) AddProfileSource(profile string, source hdevcore.PropertySource) {
	psfps.mu.Lock()
	defer psfps.mu.Unlock()
	
	psfps.profileSources[profile] = source
}