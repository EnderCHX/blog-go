package route

import (
	"blog-go/internal/interfaces/handle"
	"github.com/gin-gonic/gin"
)

type ChatRoute struct {
	ChatHandle *handle.ChatHandle
}

func NewChatRouteRegister(chatHandle *handle.ChatHandle) *ChatRoute {
	return &ChatRoute{
		ChatHandle: chatHandle,
	}
}

func (c *ChatRoute) Register(router *gin.Engine) {
	router.GET("/chat", c.ChatHandle.Chat)
}
