package route

import (
	"blog-go/internal/interfaces/handle"
	"blog-go/internal/interfaces/midware"
	"github.com/gin-gonic/gin"
)

type CommentRoute struct {
	CommentHandle *handle.CommentHandle
}

func NewCommentRouteRegister(commentHandle *handle.CommentHandle) *CommentRoute {
	return &CommentRoute{CommentHandle: commentHandle}
}
func (c CommentRoute) Register(router *gin.Engine) {
	router.GET("/comments/passage/:passageId", c.CommentHandle.GetCommentsByPassageId)
	router.GET("/comments/user/:username", c.CommentHandle.GetCommentsByUser)
	router.GET("/comment/:commentId", c.CommentHandle.GetCommentById)
	router.POST("/newcomment", midware.Permission(), c.CommentHandle.AddComment)
	router.DELETE("/comment/:commentId", c.CommentHandle.DeleteComment)
}
