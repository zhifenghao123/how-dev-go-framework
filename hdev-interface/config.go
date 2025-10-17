package hdevInterface

// IConfig 配置接口
type IConfig interface {
	// Load 加载配置
	Load(string, ...string) error
	// Get 获取配置值，返回接口类型的值
	Get(string, ...func(i interface{})) interface{}
	// GetAsString 获取配置值，返回字符串类型的值
	GetAsString(string, ...func(i interface{})) string
	// GetAsInt 获取配置值，返回整型类型的值
	GetAsInt(string, ...func(i interface{})) int
	// GetAsFloat 获取配置值，返回浮点型类型的值
	GetAsFloat(string, ...func(i interface{})) float64
	// GetAsBool 获取配置值，返回布尔类型的值
	GetAsBool(string, ...func(i interface{})) bool
	// GetAsArray 获取配置值，返回数组类型的值
	GetAsArray(string, ...func(i interface{})) []interface{}
	// GetAsMap 获取配置值，返回map类型的值
	GetAsMap(string, ...func(i interface{})) map[string]interface{}
	// GetAsStruct 获取配置值，返回结构体类型的值
	GetAsStruct(string, interface{}, ...func(i interface{})) interface{}
	// GetAsStructArray 获取配置值，返回结构体数组类型的值
	GetAsStructArray(string, interface{}, ...func(i interface{})) []interface{}
}
