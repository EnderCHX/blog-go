package midware

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/interfaces/response"
	"github.com/gin-gonic/gin"
)

func Permission() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, ok := c.Get("role")
		if !ok {
			c.JSON(response.Unauthorized(nil))
			c.Abort()
			return
		}

		role = role.(string)

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
