package hdev_ioc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	
	hdevcore "github.com/zhifenghao123/how-dev-go-framework/hdev-core"
)

// Component 组件注解
type Component struct {
	Name string
}

// Service 服务注解
type Service struct {
	Name string
}

// Repository 仓储注解
type Repository struct {
	Name string
}

// Controller 控制器注解
type Controller struct {
	Name string
}

// Configuration 配置注解
type Configuration struct{}

// Bean Bean注解
type Bean struct {
	Name string
}

// Autowired 自动注入注解
type Autowired struct{}

// Value 值注解
type Value struct {
	Value string
}

// AnnotationScanner 注解扫描器
type AnnotationScanner struct {
	basePackages []string
	beanFactory  BeanFactory
}

// NewAnnotationScanner 创建注解扫描器
func NewAnnotationScanner(beanFactory BeanFactory, basePackages ...string) *AnnotationScanner {
	return &AnnotationScanner{
		basePackages: basePackages,
		beanFactory:  beanFactory,
	}
}

// Scan 扫描包并注册Bean
func (scanner *AnnotationScanner) Scan() error {
	for _, basePackage := range scanner.basePackages {
		err := scanner.scanPackage(basePackage)
		if err != nil {
			return err
		}
	}
	return nil
}

// scanPackage 扫描包
func (scanner *AnnotationScanner) scanPackage(packagePath string) error {
	// 获取包的绝对路径
	absPath, err := scanner.findPackagePath(packagePath)
	if err != nil {
		return err
	}
	
	// 遍历包目录
	return filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 只处理.go文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			return scanner.scanFile(path)
		}
		
		return nil
	})
}

// scanFile 扫描文件
func (scanner *AnnotationScanner) scanFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	
	// 遍历文件中的声明
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			// 处理类型声明
			if x.Tok == token.TYPE {
				scanner.processTypeDecl(x, filePath)
			}
		case *ast.FuncDecl:
			// 处理方法声明（用于@Bean注解）
			scanner.processFuncDecl(x, filePath)
		}
		return true
	})
	
	return nil
}

// processTypeDecl 处理类型声明
func (scanner *AnnotationScanner) processTypeDecl(decl *ast.GenDecl, filePath string) {
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		
		// 检查是否有注解
		annotations := scanner.extractAnnotations(decl.Doc)
		
		// 注册组件
		scanner.registerComponent(typeSpec.Name.Name, annotations, filePath)
	}
}

// processFuncDecl 处理方法声明
func (scanner *AnnotationScanner) processFuncDecl(decl *ast.FuncDecl, filePath string) {
	// 检查是否有@Bean注解
	annotations := scanner.extractAnnotations(decl.Doc)
	
	for _, annotation := range annotations {
		if annotation.Name == "Bean" {
			scanner.registerBeanMethod(decl, annotation, filePath)
			break
		}
	}
}

// extractAnnotations 提取注解
func (scanner *AnnotationScanner) extractAnnotations(doc *ast.CommentGroup) []Annotation {
	if doc == nil {
		return nil
	}
	
	annotations := make([]Annotation, 0)
	
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		
		// 检查是否是注解（以@开头）
		if strings.HasPrefix(text, "@") {
			annotation := scanner.parseAnnotation(text)
			if annotation != nil {
				annotations = append(annotations, *annotation)
			}
		}
	}
	
	return annotations
}

// parseAnnotation 解析注解
func (scanner *AnnotationScanner) parseAnnotation(text string) *Annotation {
	// 移除@符号
	text = strings.TrimPrefix(text, "@")
	
	parts := strings.SplitN(text, "(", 2)
	name := strings.TrimSpace(parts[0])
	
	annotation := &Annotation{Name: name}
	
	// 解析参数
	if len(parts) > 1 {
		paramsText := strings.TrimSuffix(parts[1], ")")
		annotation.Params = scanner.parseParams(paramsText)
	}
	
	return annotation
}

// parseParams 解析参数
func (scanner *AnnotationScanner) parseParams(paramsText string) map[string]string {
	params := make(map[string]string)
	
	paramPairs := strings.Split(paramsText, ",")
	for _, pair := range paramPairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(strings.Trim(kv[1], `"`))
			params[key] = value
		}
	}
	
	return params
}

// registerComponent 注册组件
func (scanner *AnnotationScanner) registerComponent(typeName string, annotations []Annotation, filePath string) {
	var beanName string
	var scope BeanScope = ScopeSingleton
	
	// 检查注解类型
	for _, annotation := range annotations {
		switch annotation.Name {
		case "Component", "Service", "Repository", "Controller":
			if name, exists := annotation.Params["name"]; exists {
				beanName = name
			} else {
				beanName = strings.ToLower(typeName[:1]) + typeName[1:]
			}
		case "Scope":
			if scopeValue, exists := annotation.Params["value"]; exists {
				scope = BeanScope(scopeValue)
			}
		}
	}
	
	if beanName == "" {
		return // 没有找到组件注解
	}
	
	// 注册Bean定义
	beanType := scanner.resolveType(typeName, filePath)
	if beanType == nil {
		return
	}
	
	beanDefinition := NewBeanDefinition(beanName, beanType)
	beanDefinition.SetScope(scope)
	
	err := scanner.beanFactory.RegisterBeanDefinition(beanName, beanDefinition)
	if err != nil {
		fmt.Printf("Failed to register bean %s: %v\n", beanName, err)
	}
}

// registerBeanMethod 注册Bean方法
func (scanner *AnnotationScanner) registerBeanMethod(decl *ast.FuncDecl, annotation Annotation, filePath string) {
	// 简化实现
	fmt.Printf("Found @Bean method: %s\n", decl.Name.Name)
}

// resolveType 解析类型
func (scanner *AnnotationScanner) resolveType(typeName string, filePath string) reflect.Type {
	// 简化实现：使用反射查找类型
	// 实际实现应该更复杂，需要考虑包导入等
	return nil
}

// findPackagePath 查找包路径
func (scanner *AnnotationScanner) findPackagePath(packagePath string) (string, error) {
	// 简化实现：在当前工作目录下查找
	return packagePath, nil
}

// Annotation 注解结构
type Annotation struct {
	Name   string
	Params map[string]string
}

// ClassPathScanner 类路径扫描器
type ClassPathScanner struct {
	basePackages []string
}

// NewClassPathScanner 创建类路径扫描器
func NewClassPathScanner(basePackages ...string) *ClassPathScanner {
	return &ClassPathScanner{
		basePackages: basePackages,
	}
}

// Scan 扫描类路径
func (scanner *ClassPathScanner) Scan() ([]string, error) {
	classes := make([]string, 0)
	
	for _, basePackage := range scanner.basePackages {
		packageClasses, err := scanner.scanPackage(basePackage)
		if err != nil {
			return nil, err
		}
		classes = append(classes, packageClasses...)
	}
	
	return classes, nil
}

// scanPackage 扫描包
func (scanner *ClassPathScanner) scanPackage(packagePath string) ([]string, error) {
	classes := make([]string, 0)
	
	// 简化实现
	return classes, nil
}

// ComponentScan 组件扫描配置
type ComponentScan struct {
	BasePackages []string
}

// EnableComponentScan 启用组件扫描
// 注意：调用方需要自行提供本模块的 BeanFactory 实例进行扫描
func EnableComponentScan(beanFactory BeanFactory, basePackages ...string) hdevcore.ApplicationContextInitializer {
	return func(ctx hdevcore.ApplicationContext) error {
		if beanFactory == nil {
			return fmt.Errorf("beanFactory must not be nil")
		}
		scanner := NewAnnotationScanner(beanFactory, basePackages...)
		return scanner.Scan()
	}
}

// AutoConfiguration 自动配置注解
type AutoConfiguration struct{}

// Conditional 条件注解
type Conditional struct {
	Condition string
}

// ConditionalOnClass 类存在条件
type ConditionalOnClass struct {
	Classes []string
}

// ConditionalOnMissingClass 类不存在条件
type ConditionalOnMissingClass struct {
	Classes []string
}

// ConditionalOnBean Bean存在条件
type ConditionalOnBean struct {
	Beans []string
}

// ConditionalOnMissingBean Bean不存在条件
type ConditionalOnMissingBean struct {
	Beans []string
}

// ConditionalOnProperty 属性条件
type ConditionalOnProperty struct {
	Name     string
	HavingValue string
	MatchIfMissing bool
}

// Example usage:
/*
// UserService.go
package service

// @Service(name="userService")
type UserService struct {
	// @Autowired
	UserRepository *UserRepository
}

func (s *UserService) AfterPropertiesSet() error {
	// 初始化逻辑
	return nil
}

// Application.go
package main

import (
	"github.com/zhifenghao123/how-dev-go-framework/hdev-ioc"
)

func main() {
	// 创建Bean工厂
	beanFactory := hdevioc.NewDefaultBeanFactory()
	
	// 启用组件扫描
	scanner := hdevioc.NewAnnotationScanner(beanFactory, "com.example.service")
	scanner.Scan()
	
	// 预实例化单例Bean
	beanFactory.PreInstantiateSingletons()
	
	// 获取Bean实例
	userService := beanFactory.Instance(&service.UserService{}).(*service.UserService)
	
	fmt.Printf("UserService: %+v\n", userService)
}
*/