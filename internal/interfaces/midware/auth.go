package midware

import (
	"blog-go/internal/infrastructure/config"
	"net/http"
	"strings"

	"github.com/EnderCHX/chx-tools-go/auth"
	"github.com/gin-gonic/gin"
)

func Auth(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.Replace(token, "Bearer ", "", 1)
			claims, err := auth.VerifyToken(token, config.SecretKeys.AccessSecret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "token非法",
					"code":    "InvalidAccessToken",
					"data":    nil,
				})
				c.Abort()
			} else {
				c.Set("claims", claims)
				c.Next()
			}
		}
	}
}
