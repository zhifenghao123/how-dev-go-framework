package hdevgo

import (
	"log"
	"reflect"
	"runtime/debug"
	"sync"
)

// Process 进程模型，对应 yaml 中 process: 配置项的每一个条目
//
//	process:
//	  - class: HttpServer
//	    execute: Execute
//	    params:
//	      Ip: 0.0.0.0
//	      Port: 8080
type Process struct {
	Id      int
	Name    string
	Class   string
	Execute string
	Params  map[string]interface{}
}

// 配置 key 常量
const (
	configKeyProcess        = "process"
	configKeyImport         = "import"
	defaultProcessExecuteFn = "Execute"
)

// extractProcesses 从配置中读取 process 列表
// cname 为顶层配置文件名（不带扩展名），用于拼接完整 key（兼容旧版 hdev-go 使用 "{cname}.process"）
// 若按 "{cname}.process" 取不到则尝试直接 "process"
func extractProcesses(cfg *ConfigPropertySource, cname string) []*Process {
	if cfg == nil {
		return nil
	}
	rawList := cfg.cfg.GetAsStructArray(joinDot(cname, configKeyProcess), Process{})
	if rawList == nil || len(rawList) == 0 {
		rawList = cfg.cfg.GetAsStructArray(configKeyProcess, Process{})
	}
	if rawList == nil {
		return nil
	}
	processes := make([]*Process, 0, len(rawList))
	for _, item := range rawList {
		if p, ok := item.(*Process); ok {
			processes = append(processes, p)
		}
	}
	return processes
}

// runProcess 调度单个 Process：从容器取出 Bean，注入 Params，调用 Execute 方法
// wg 用于通知调度器该 Process 已结束
func runProcess(p *Process, beanLookup func(name string) interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[hdev-context] runProcess panic: class=%s err=%v", p.Class, r)
			debug.PrintStack()
		}
	}()

	bean := beanLookup(p.Class)
	if bean == nil {
		log.Printf("[hdev-context] runProcess: bean %s not found", p.Class)
		return
	}

	vl := reflect.ValueOf(bean)
	// 把 Params 中的字段反射注入
	if vl.Kind() == reflect.Ptr && !vl.IsNil() {
		elem := vl.Elem()
		for k, v := range p.Params {
			field := elem.FieldByName(k)
			if !field.IsValid() {
				log.Printf("[hdev-context] runProcess: field %s not found in %s", k, p.Class)
				continue
			}
			if !field.CanSet() {
				log.Printf("[hdev-context] runProcess: field %s in %s cannot set", k, p.Class)
				continue
			}
			setProcessField(field, reflect.ValueOf(v))
		}
	}

	execName := p.Execute
	if execName == "" {
		execName = defaultProcessExecuteFn
	}
	exec := vl.MethodByName(execName)
	if !exec.IsValid() {
		log.Printf("[hdev-context] runProcess: method %s not found on %s", execName, p.Class)
		return
	}
	exec.Call(nil)
}

// setProcessField 反射安全地为目标字段赋值（兼容类型转换）
func setProcessField(f reflect.Value, v reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[hdev-context] inject process param error: %v", err)
		}
	}()

	if !v.IsValid() {
		return
	}

	if f.Kind() == v.Kind() || (f.Kind() == reflect.Interface && v.Type().Implements(f.Type())) {
		f.Set(v)
		return
	}

	// 兼容数值/字符串等可转换类型
	if v.Type().ConvertibleTo(f.Type()) {
		f.Set(v.Convert(f.Type()))
		return
	}

	log.Printf("[hdev-context] cannot inject %s with %s", f.Kind().String(), v.Kind().String())
}
