package midware

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/interfaces/response"
	"github.com/EnderCHX/chx-tools-go/auth"
	"github.com/gin-gonic/gin"
)

func Permission() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")

		role := claims.(*auth.JWTPayload).Role

		if c.FullPath() == "/newpost" &&
			(role == entity.ADMIN || role == entity.EDITOR || role == entity.AUTHOR) {
			c.Next()
		} else if c.FullPath() == "/newcomment" &&
			(role == entity.ADMIN || role == entity.EDITOR || role == entity.AUTHOR || role == entity.USER) {
			c.Next()
		} else {
			c.JSON(response.Default(200, false, response.USER_PERMISSION_DENIED, "权限不足", nil))
			c.Abort()
		}
	}
}
