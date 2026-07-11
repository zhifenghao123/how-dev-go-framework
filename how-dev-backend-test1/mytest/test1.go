package mytest

import (
	"fmt"
)

// HTTPServer 结构体
type HTTPServer struct {
	addr    string
	timeout int
}

// NewHTTPServer 创建HTTPServer实例的构造函数
func NewHTTPServer(options ...func(*HTTPServer)) *HTTPServer {
	server := &HTTPServer{
		addr:    ":8080", // 默认地址
		timeout: 10,      // 默认超时时间
	}

	// 应用所有选项
	for _, option := range options {
		option(server)
	}

	return server
}

// WithAddr 设置地址的选项函数
func WithAddr(addr string) func(*HTTPServer) {
	return func(s *HTTPServer) {
		s.addr = addr
	}
}

// WithTimeout 设置超时时间的选项函数
func WithTimeout(timeout int) func(*HTTPServer) {
	return func(s *HTTPServer) {
		s.timeout = timeout
	}
}

func main() {
	server := NewHTTPServer(WithAddr(":9090"), WithTimeout(20))
	fmt.Printf("Server Address: %s, Timeout: %d\n", server.addr, server.timeout)
}
