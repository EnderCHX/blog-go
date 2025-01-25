package midware

import (
	"blog-go/config"
	"net/http"
	"strings"

	"github.com/EnderCHX/chx-tools-go/auth"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.Replace(token, "Bearer ", "", 1)
			claims, err := auth.VerifyToken(token, config.ConfigContext.SecretKeys.AccessSecret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": err,
					"code":    "InvalidAccessToken",
					"data":    nil,
				})
				c.Abort()
				return
			} else {
				c.Set("claims", claims)
				c.Next()
			}
		}
	}
}
