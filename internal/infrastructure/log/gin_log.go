package log

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		Logger.Info("[GIN]HTTP 请求",
			zap.String(" status", fmt.Sprintf("%v", c.Writer.Status())),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("latency", fmt.Sprintf("%v", latency)),
		)
	}
}
