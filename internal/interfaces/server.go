package interfaces

import (
	"blog-go/internal/interfaces/route"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Start()
}

type HttpServer struct {
	router  *gin.Engine
	routes  []route.RouteRegister
	midware []gin.HandlerFunc
	ip      string
	port    string
	level   string
}

func NewHttpServer(ip, port, level string, middleware []gin.HandlerFunc, routes ...route.RouteRegister) *HttpServer {
	gin.SetMode(level)
	r := gin.New()
	return &HttpServer{
		router:  r,
		ip:      ip,
		port:    port,
		level:   level,
		midware: middleware,
		routes:  routes,
	}
}

func (s *HttpServer) Start() {
	s.router.Use(s.midware...)

	for _, r := range s.routes {
		r.Register(s.router)
	}

	s.router.Run(s.ip + ":" + s.port)
}
