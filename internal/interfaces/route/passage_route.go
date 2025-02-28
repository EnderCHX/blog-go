package route

import (
	"blog-go/internal/interfaces/handle"

	"github.com/gin-gonic/gin"
)

type PassageRouteRegister struct {
	passageHandle *handle.PassageHandle
}

func NewPassageRouteRegister(passageHandle *handle.PassageHandle) *PassageRouteRegister {
	return &PassageRouteRegister{
		passageHandle: passageHandle,
	}
}

func (r *PassageRouteRegister) Register(router *gin.Engine) {
	router.GET("/posts", r.passageHandle.GetPassages)
	router.GET("/posts/:tagname", r.passageHandle.GetPassagesByTag)
	// r.router.GET("/posts/tags", r.passageHandle.GetTags)
	router.GET("/posts/tags/:tagname", r.passageHandle.GetPassagesByTag)
	router.GET("/post/:id", r.passageHandle.GetPassage)
	router.POST("/newpost", r.passageHandle.AddPassage)
}
