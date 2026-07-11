package hdevcore

import (
	"io"
)

// Resource 资源接口
type Resource interface {
	// Exists 检查资源是否存在
	Exists() bool
	
	// IsReadable 检查资源是否可读
	IsReadable() bool
	
	// IsOpen 检查资源是否已打开
	IsOpen() bool
	
	// GetURL 获取资源的URL
	GetURL() (string, error)
	
	// GetFile 获取资源的文件
	GetFile() (interface{}, error)
	
	// ContentLength 获取内容长度
	ContentLength() (int64, error)
	
	// LastModified 获取最后修改时间
	LastModified() (int64, error)
	
	// GetFilename 获取文件名
	GetFilename() string
	
	// GetDescription 获取资源描述
	GetDescription() string
	
	// GetInputStream 获取输入流
	GetInputStream() (io.Reader, error)
}

// WritableResource 可写资源接口
type WritableResource interface {
	Resource
	
	// IsWritable 检查资源是否可写
	IsWritable() bool
	
	// GetOutputStream 获取输出流
	GetOutputStream() (io.Writer, error)
	
	// Write 写入内容
	Write(data []byte) error
}

// ContextResource 上下文资源接口
type ContextResource interface {
	Resource
	
	// GetPathWithinContext 获取上下文内的路径
	GetPathWithinContext() string
}

// ResourceLoader 资源加载器接口
type ResourceLoader interface {
	// GetResource 获取资源
	GetResource(location string) (Resource, error)
	
	// GetClassLoader 获取类加载器
	GetClassLoader() interface{}
}

// ResourcePatternResolver 资源模式解析器接口
type ResourcePatternResolver interface {
	ResourceLoader
	
	// GetResources 获取匹配模式的所有资源
	GetResources(locationPattern string) ([]Resource, error)
}

// ResourceResolver 资源解析器接口
type ResourceResolver interface {
	// ResolveResource 解析资源
	ResolveResource(location string) (Resource, error)
}

// ResourcePattern 资源模式接口
type ResourcePattern interface {
	// Matches 检查是否匹配
	Matches(resource Resource) bool
}