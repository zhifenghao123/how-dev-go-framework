package hdevcore

import (
	"time"
)

// ApplicationEvent 应用事件接口
type ApplicationEvent interface {
	// GetTimestamp 获取事件时间戳
	GetTimestamp() time.Time
	
	// GetSource 获取事件源
	GetSource() interface{}
}

// ApplicationListener 应用事件监听器接口
type ApplicationListener interface {
	// OnApplicationEvent 处理应用事件
	OnApplicationEvent(event ApplicationEvent)
}

// ApplicationEventPublisher 应用事件发布器接口
type ApplicationEventPublisher interface {
	// PublishEvent 发布事件
	PublishEvent(event ApplicationEvent)
}

// ApplicationEventMulticaster 应用事件广播器接口
type ApplicationEventMulticaster interface {
	ApplicationEventPublisher
	
	// AddApplicationListener 添加应用监听器
	AddApplicationListener(listener ApplicationListener)
	
	// AddApplicationListenerBean 通过Bean名称添加应用监听器
	AddApplicationListenerBean(beanName string)
	
	// RemoveApplicationListener 移除应用监听器
	RemoveApplicationListener(listener ApplicationListener)
	
	// RemoveAllListeners 移除所有监听器
	RemoveAllListeners()
}

// SmartApplicationListener 智能应用监听器接口
type SmartApplicationListener interface {
	ApplicationListener
	
	// SupportsEventType 支持的事件类型
	SupportsEventType(eventType interface{}) bool
	
	// SupportsSourceType 支持的源类型
	SupportsSourceType(sourceType interface{}) bool
}

// GenericApplicationEvent 通用应用事件
type GenericApplicationEvent struct {
	timestamp time.Time
	source    interface{}
}

// NewGenericApplicationEvent 创建通用应用事件
func NewGenericApplicationEvent(source interface{}) *GenericApplicationEvent {
	return &GenericApplicationEvent{
		timestamp: time.Now(),
		source:    source,
	}
}

func (e *GenericApplicationEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e *GenericApplicationEvent) GetSource() interface{} {
	return e.source
}