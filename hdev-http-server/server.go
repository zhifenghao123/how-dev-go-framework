package ihttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	DefaultIp           = "127.0.0.1"
	DefaultPort         = 80
	DefaultReadTimeout  = 1
	DefaultWriteTimeout = 30
	DefaultHandlerMode  = 0
)

// HttpServer http server
type HttpServer struct {
	*http.Server
	Ctx          context.Context `wired:"true"`
	Ioc          IBeanFactory    `wired:"true"`
	Name         string
	Ip           string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Route        string
	HandlerMode  int
	router       *RouterConstruct
	handler      *HandlerConstruct
}

// Execute execute function to start server
func (s *HttpServer) Execute() {
	s.options()

	h := NewHandler(s.HandlerMode, s.handler)
	s.Server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Ip, s.Port),
		ReadTimeout:  time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.WriteTimeout) * time.Second,
		Handler: h.Handle(
			NewRouter(s.router).Action(),
		),
	}

	go watch(s, h)

	log.Printf("=== 【Http】Server [%s] start [%s] ===", s.Name, s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Printf("http server [%s] stop [%s]", s.Name, err)
	}
}

// SetRouter set http router
func (s *HttpServer) SetRouter(r IRouter, opt ...IRouterOption) {
	if r == nil {
		return
	}
	s.router = &RouterConstruct{
		Router:  r,
		Options: opt,
	}
}

// SetHandler set request handler
func (s *HttpServer) SetHandler(h IHandler, opt ...IHandlerOption) {
	if h == nil {
		return
	}
	s.handler = &HandlerConstruct{
		Handler: h,
		Options: opt,
	}
}

func (s *HttpServer) options() {
	if s.Ip == "" {
		s.Ip = DefaultIp
	}
	if s.Port == 0 {
		s.Port = DefaultPort
	}
	if s.ReadTimeout == 0 {
		s.ReadTimeout = DefaultReadTimeout
	}
	if s.WriteTimeout == 0 {
		s.WriteTimeout = DefaultWriteTimeout
	}
	if s.HandlerMode != 0 && s.HandlerMode != 1 {
		// 不支持的处理模式，则异常终止启动
		panic(fmt.Sprintf("unsupported handler mode [%d], only support 0 or 1", s.HandlerMode))
	}

	if s.router == nil {
		s.router = newRouter(WithRouterFile(s.Route))
	}
	if s.handler == nil {
		s.handler = newHandler(s.HandlerMode, WithIoc(s.Ioc))
	}
}

func watch(s *HttpServer, h IHandler) {
	if s.Ctx == nil {
		return
	}

	for {
		select {
		case <-s.Ctx.Done():
			log.Println("http server stop ...")
			log.Println("wait for all http requests return ...")
			h.Done()
			err := s.Shutdown(s.Ctx)
			if err != nil {
				log.Println("http server shutdown error : ", err)
			}
			log.Println("http server stop success")
			return
		}
	}
}
