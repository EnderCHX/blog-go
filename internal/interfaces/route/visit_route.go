package route

import (
	"blog-go/internal/interfaces/handle"
	"github.com/gin-gonic/gin"
)

type VisitRoute struct {
	VisitHandle *handle.VisitHandle
}

func NewVisitRoute(visitHandle *handle.VisitHandle) *VisitRoute {
	return &VisitRoute{
		VisitHandle: visitHandle,
	}
}

func (r *VisitRoute) Register(router *gin.Engine) {
	router.GET("/visit", r.VisitHandle.GetVisitCount)
}
