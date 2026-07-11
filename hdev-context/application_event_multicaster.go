package hdevcontext

import (
	"reflect"
	"sync"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// SimpleApplicationEventMulticaster 简单应用事件广播器实现
type SimpleApplicationEventMulticaster struct {
	listeners      []hdevcore.ApplicationListener
	listenerBeans  []string
	defaultExecutor interface{}
	mu             sync.RWMutex
}

// NewSimpleApplicationEventMulticaster 创建简单事件广播器
func NewSimpleApplicationEventMulticaster() *SimpleApplicationEventMulticaster {
	return &SimpleApplicationEventMulticaster{
		listeners:     make([]hdevcore.ApplicationListener, 0),
		listenerBeans: make([]string, 0),
	}
}

// PublishEvent 发布事件
func (multicaster *SimpleApplicationEventMulticaster) PublishEvent(event hdevcore.ApplicationEvent) {
	multicaster.mu.RLock()
	defer multicaster.mu.RUnlock()
	
	// 同步通知所有监听器
	for _, listener := range multicaster.listeners {
		multicaster.invokeListener(listener, event)
	}
}

// AddApplicationListener 添加应用监听器
func (multicaster *SimpleApplicationEventMulticaster) AddApplicationListener(listener hdevcore.ApplicationListener) {
	multicaster.mu.Lock()
	defer multicaster.mu.Unlock()
	
	multicaster.listeners = append(multicaster.listeners, listener)
}

// AddApplicationListenerBean 通过Bean名称添加应用监听器
func (multicaster *SimpleApplicationEventMulticaster) AddApplicationListenerBean(beanName string) {
	multicaster.mu.Lock()
	defer multicaster.mu.Unlock()
	
	multicaster.listenerBeans = append(multicaster.listenerBeans, beanName)
}

// RemoveApplicationListener 移除应用监听器
func (multicaster *SimpleApplicationEventMulticaster) RemoveApplicationListener(listener hdevcore.ApplicationListener) {
	multicaster.mu.Lock()
	defer multicaster.mu.Unlock()
	
	for i, l := range multicaster.listeners {
		if l == listener {
			multicaster.listeners = append(multicaster.listeners[:i], multicaster.listeners[i+1:]...)
			break
		}
	}
}

// RemoveAllListeners 移除所有监听器
func (multicaster *SimpleApplicationEventMulticaster) RemoveAllListeners() {
	multicaster.mu.Lock()
	defer multicaster.mu.Unlock()
	
	multicaster.listeners = make([]hdevcore.ApplicationListener, 0)
	multicaster.listenerBeans = make([]string, 0)
}

// invokeListener 调用监听器
func (multicaster *SimpleApplicationEventMulticaster) invokeListener(listener hdevcore.ApplicationListener, event hdevcore.ApplicationEvent) {
	// 如果是智能监听器，检查是否支持该事件
	if smartListener, ok := listener.(hdevcore.SmartApplicationListener); ok {
		if !smartListener.SupportsEventType(reflect.TypeOf(event)) {
			return
		}
		if !smartListener.SupportsSourceType(reflect.TypeOf(event.GetSource())) {
			return
		}
	}
	
	// 调用监听器
	listener.OnApplicationEvent(event)
}

// GenericApplicationListenerAdapter 通用应用监听器适配器
type GenericApplicationListenerAdapter struct {
	delegate interface{}
}

// NewGenericApplicationListenerAdapter 创建通用监听器适配器
func NewGenericApplicationListenerAdapter(delegate interface{}) *GenericApplicationListenerAdapter {
	return &GenericApplicationListenerAdapter{
		delegate: delegate,
	}
}

func (adapter *GenericApplicationListenerAdapter) OnApplicationEvent(event hdevcore.ApplicationEvent) {
	// 这里应该通过反射调用delegate的方法
	// 简化实现
}

// ApplicationEventMulticasterSupport 应用事件广播器支持
type ApplicationEventMulticasterSupport struct {
	multicaster hdevcore.ApplicationEventMulticaster
}

// NewApplicationEventMulticasterSupport 创建事件广播器支持
func NewApplicationEventMulticasterSupport() *ApplicationEventMulticasterSupport {
	return &ApplicationEventMulticasterSupport{
		multicaster: NewSimpleApplicationEventMulticaster(),
	}
}

// GetApplicationEventMulticaster 获取事件广播器
func (support *ApplicationEventMulticasterSupport) GetApplicationEventMulticaster() hdevcore.ApplicationEventMulticaster {
	return support.multicaster
}

// SetApplicationEventMulticaster 设置事件广播器
func (support *ApplicationEventMulticasterSupport) SetApplicationEventMulticaster(multicaster hdevcore.ApplicationEventMulticaster) {
	support.multicaster = multicaster
}

// PublishEvent 发布事件
func (support *ApplicationEventMulticasterSupport) PublishEvent(event hdevcore.ApplicationEvent) {
	if support.multicaster != nil {
		support.multicaster.PublishEvent(event)
	}
}

// AddApplicationListener 添加应用监听器
func (support *ApplicationEventMulticasterSupport) AddApplicationListener(listener hdevcore.ApplicationListener) {
	if support.multicaster != nil {
		support.multicaster.AddApplicationListener(listener)
	}
}

// AddApplicationListenerBean 通过Bean名称添加应用监听器
func (support *ApplicationEventMulticasterSupport) AddApplicationListenerBean(beanName string) {
	if support.multicaster != nil {
		support.multicaster.AddApplicationListenerBean(beanName)
	}
}

// ApplicationListenerMethodAdapter 应用监听器方法适配器
type ApplicationListenerMethodAdapter struct {
	beanName     string
	methodName   string
	targetObject interface{}
}

// NewApplicationListenerMethodAdapter 创建监听器方法适配器
func NewApplicationListenerMethodAdapter(beanName, methodName string, targetObject interface{}) *ApplicationListenerMethodAdapter {
	return &ApplicationListenerMethodAdapter{
		beanName:     beanName,
		methodName:   methodName,
		targetObject: targetObject,
	}
}

func (adapter *ApplicationListenerMethodAdapter) OnApplicationEvent(event hdevcore.ApplicationEvent) {
	// 这里应该通过反射调用目标方法
	// 简化实现
}

// PayloadApplicationEvent 负载应用事件
type PayloadApplicationEvent struct {
	*hdevcore.GenericApplicationEvent
	payload interface{}
}

// NewPayloadApplicationEvent 创建负载应用事件
func NewPayloadApplicationEvent(source, payload interface{}) *PayloadApplicationEvent {
	return &PayloadApplicationEvent{
		GenericApplicationEvent: hdevcore.NewGenericApplicationEvent(source),
		payload:                payload,
	}
}

// GetPayload 获取负载
func (event *PayloadApplicationEvent) GetPayload() interface{} {
	return event.payload
}

// OrderedApplicationListener 有序应用监听器包装器
type OrderedApplicationListener struct {
	delegate hdevcore.ApplicationListener
	order    int
}

// NewOrderedApplicationListener 创建有序应用监听器
func NewOrderedApplicationListener(delegate hdevcore.ApplicationListener, order int) *OrderedApplicationListener {
	return &OrderedApplicationListener{
		delegate: delegate,
		order:    order,
	}
}

func (listener *OrderedApplicationListener) OnApplicationEvent(event hdevcore.ApplicationEvent) {
	listener.delegate.OnApplicationEvent(event)
}

// GetOrder 获取排序值
func (listener *OrderedApplicationListener) GetOrder() int {
	return listener.order
}