package hdevconfig

// IProcessor 配置处理器接口
type IProcessor interface {
	LoadFileContentAsMapData(file string) (map[string]interface{}, error)
}

// ProcessorConstruct 用于构造新IProcessor实例的函数类型
type ProcessorConstruct func() IProcessor
