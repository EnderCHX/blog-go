package route

import (
	"blog-go/internal/interfaces/handle"
	"github.com/gin-gonic/gin"
)

type TagRouteRegister struct {
	tagHandle *handle.TagHandle
}

func NewTagRouteRegister(tagHandle *handle.TagHandle) *TagRouteRegister {
	return &TagRouteRegister{
		tagHandle: tagHandle,
	}
}

func (r *TagRouteRegister) Register(router *gin.Engine) {
	router.GET("/posts/tags", r.tagHandle.GetTags)
	router.GET("/post/:id/tags", r.tagHandle.GetPassageTags)
}
