package hdevconfig

// IExtractor 配置提取器接口
type IExtractor interface {
	Load(interface{}, string)
	Get(string, ...func(i interface{})) interface{}
	GetAsString(string, ...func(i interface{})) string
	GetAsInt(string, ...func(i interface{})) int
	GetAsFloat(string, ...func(i interface{})) float64
	GetAsBool(string, ...func(i interface{})) bool
	GetAsArray(string, ...func(i interface{})) []interface{}
	GetAsMap(string, ...func(i interface{})) map[string]interface{}
	GetAsStruct(string, interface{}, ...func(i interface{})) interface{}
	GetAsStructArray(string, interface{}, ...func(i interface{})) []interface{}
}

// ExtractorConstruct 用于构造新IExtractor实例的函数类型
type ExtractorConstruct func() IExtractor
