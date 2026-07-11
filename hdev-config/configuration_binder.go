package hdevconfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// ConfigurationBinder 配置绑定器
type ConfigurationBinder struct {
	environment hdevcore.Environment
	validators  []ConfigurationValidator
	mu          sync.RWMutex
}

// NewConfigurationBinder 创建配置绑定器
func NewConfigurationBinder(environment hdevcore.Environment) *ConfigurationBinder {
	return &ConfigurationBinder{
		environment: environment,
		validators:  make([]ConfigurationValidator, 0),
	}
}

// Bind 绑定配置到目标对象
func (binder *ConfigurationBinder) Bind(prefix string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	
	targetElem := targetValue.Elem()
	if targetElem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}
	
	return binder.bindStruct(prefix, targetElem)
}

// BindYaml 从YAML数据绑定配置
func (binder *ConfigurationBinder) BindYaml(data []byte, target interface{}) error {
	// 简化实现，实际应该解析YAML数据
	return nil
}

// BindProperties 从Properties数据绑定配置
func (binder *ConfigurationBinder) BindProperties(data []byte, target interface{}) error {
	// 简化实现，实际应该解析Properties数据
	return nil
}

// AddValidator 添加配置验证器
func (binder *ConfigurationBinder) AddValidator(validator ConfigurationValidator) {
	binder.mu.Lock()
	defer binder.mu.Unlock()
	
	binder.validators = append(binder.validators, validator)
}

// bindStruct 绑定结构体
func (binder *ConfigurationBinder) bindStruct(prefix string, structValue reflect.Value) error {
	structType := structValue.Type()
	
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		
		// 跳过不可导出的字段
		if !fieldValue.CanSet() {
			continue
		}
		
		// 获取字段的配置键名
		configKey := binder.getConfigKey(prefix, field)
		
		// 绑定字段值
		err := binder.bindField(configKey, field, fieldValue)
		if err != nil {
			return fmt.Errorf("failed to bind field %s: %v", field.Name, err)
		}
	}
	
	// 执行验证
	return binder.validate(structValue.Interface())
}

// bindField 绑定字段值
func (binder *ConfigurationBinder) bindField(configKey string, field reflect.StructField, fieldValue reflect.Value) error {
	// 获取配置值
	configValue := binder.environment.GetProperty(configKey)
	
	// 如果配置值为空且有默认值标签，使用默认值
	if configValue == "" {
		if defaultValue := field.Tag.Get("default"); defaultValue != "" {
			configValue = defaultValue
		}
	}
	
	// 如果配置值仍然为空且字段不是必须的，跳过
	if configValue == "" {
		if required := field.Tag.Get("required"); required == "true" {
			return fmt.Errorf("required configuration property '%s' is missing", configKey)
		}
		return nil
	}
	
	// 根据字段类型转换配置值
	convertedValue, err := binder.convertValue(configValue, fieldValue.Type())
	if err != nil {
		return fmt.Errorf("failed to convert value for field %s: %v", field.Name, err)
	}
	
	fieldValue.Set(reflect.ValueOf(convertedValue))
	return nil
}

// convertValue 转换配置值
func (binder *ConfigurationBinder) convertValue(value string, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.String:
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(v).Convert(targetType).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(v).Convert(targetType).Interface(), nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(v).Convert(targetType).Interface(), nil
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.Slice:
		return binder.convertSlice(value, targetType)
	default:
		return nil, fmt.Errorf("unsupported type: %s", targetType.Kind())
	}
}

// convertSlice 转换切片值
func (binder *ConfigurationBinder) convertSlice(value string, sliceType reflect.Type) (interface{}, error) {
	elementType := sliceType.Elem()
	values := strings.Split(value, ",")
	
	sliceValue := reflect.MakeSlice(sliceType, len(values), len(values))
	
	for i, v := range values {
		converted, err := binder.convertValue(strings.TrimSpace(v), elementType)
		if err != nil {
			return nil, err
		}
		sliceValue.Index(i).Set(reflect.ValueOf(converted))
	}
	
	return sliceValue.Interface(), nil
}

// getConfigKey 获取配置键名
// 支持的tag优先级（从高到低）：config > yaml > json > mapstructure > 字段名(snake_case)
func (binder *ConfigurationBinder) getConfigKey(prefix string, field reflect.StructField) string {
	// 检查是否有自定义键名
	var key string
	for _, tagName := range []string{"config", "yaml", "json", "mapstructure"} {
		if v := field.Tag.Get(tagName); v != "" {
			// json/yaml tag可能含有",omitempty"等修饰，取第一段
			parts := strings.Split(v, ",")
			if parts[0] != "" && parts[0] != "-" {
				key = parts[0]
				break
			}
		}
	}
	
	if key == "" {
		// 使用字段名作为键名（snake_case）
		key = binder.toSnakeCase(field.Name)
	}
	
	if prefix != "" {
		return prefix + "." + key
	}
	return key
}

// toSnakeCase 转换为蛇形命名
func (binder *ConfigurationBinder) toSnakeCase(str string) string {
	var result strings.Builder
	
	for i, char := range str {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(char)
	}
	
	return strings.ToLower(result.String())
}

// validate 验证配置
func (binder *ConfigurationBinder) validate(target interface{}) error {
	binder.mu.RLock()
	defer binder.mu.RUnlock()
	
	for _, validator := range binder.validators {
		if err := validator.Validate(target); err != nil {
			return err
		}
	}
	
	return nil
}

// ConfigurationValidator 配置验证器接口
type ConfigurationValidator interface {
	Validate(target interface{}) error
}

// RequiredFieldValidator 必需字段验证器
type RequiredFieldValidator struct{}

func (v *RequiredFieldValidator) Validate(target interface{}) error {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	
	if value.Kind() != reflect.Struct {
		return nil
	}
	
	structType := value.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)
		
		// 检查必需字段
		if required := field.Tag.Get("required"); required == "true" {
			if fieldValue.IsZero() {
				return fmt.Errorf("field %s is required but empty", field.Name)
			}
		}
	}
	
	return nil
}

// RangeValidator 范围验证器
type RangeValidator struct{}

func (v *RangeValidator) Validate(target interface{}) error {
	// 简化实现
	return nil
}

// ConfigurationProperties 配置属性注解支持
type ConfigurationProperties struct {
	Prefix string
}

// BindConfigurationProperties 绑定配置属性
func BindConfigurationProperties(environment hdevcore.Environment, prefix string, target interface{}) error {
	binder := NewConfigurationBinder(environment)
	return binder.Bind(prefix, target)
}

// ConfigurationPropertiesBinder 配置属性绑定器工厂
func ConfigurationPropertiesBinder(environment hdevcore.Environment) *ConfigurationBinder {
	return NewConfigurationBinder(environment)
}

// Example usage:
/*
type ServerConfig struct {
	Port     int    `config:"port" default:"8080"`
	Host     string `config:"host" default:"localhost"`
	SSL      bool   `config:"ssl" default:"false"`
	MaxConns int    `config:"max-connections" required:"true"`
}

func main() {
	env := NewConfigurableEnvironment()
	
	config := &ServerConfig{}
	err := BindConfigurationProperties(env, "server", config)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Server config: %+v\n", config)
}
*/