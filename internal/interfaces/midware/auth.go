package midware

import (
	"blog-go/internal/infrastructure/config"
	"blog-go/internal/interfaces/response"
	"github.com/EnderCHX/chx-tools-go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		if token != "" {
			token = strings.Replace(token, "Bearer ", "", 1)
			claims, err := auth.VerifyToken(token, config.SecretKeys.AccessSecret)
			if err != nil {
				c.JSON(response.Default(http.StatusUnauthorized, false, response.CLIENT_INVALID_ACCESS_TOKEN, "token非法", nil))
				c.Abort()
			} else {
				c.Set("claims", claims)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Next()
			}
		}
		c.Next()
	}
}
