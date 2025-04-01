package route

import (
	"github.com/gin-gonic/gin"
)

type RouteRegister interface {
	Register(router *gin.Engine)
}
