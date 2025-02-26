package interfaces

import "github.com/gin-gonic/gin"

type Server interface {
	Start()
}

type HttpServer struct {
	router *gin.Engine
	ip     string
	port   string
	level  string
}

func NewHttpServer(ip, port, level string, middleware ...gin.HandlerFunc) *HttpServer {
	gin.SetMode(level)
	r := gin.New()
	r.Use(middleware...)
	return &HttpServer{
		router: r,
		ip:     ip,
		port:   port,
		level:  level,
	}
}

func (s *HttpServer) AddRouter(method, path string, handler ...gin.HandlerFunc) *HttpServer {
	s.router.Handle(method, path, handler...)
	return s
}

func (s *HttpServer) GetRouter() *gin.Engine {
	return s.router
}

func (s *HttpServer) Start() *HttpServer {
	s.router.Run(s.ip + ":" + s.port)
	return s
}
