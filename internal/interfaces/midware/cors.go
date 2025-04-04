package midware

import (
	"blog-go/internal/infrastructure/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(cf config.Config) gin.HandlerFunc { //跨域中间件
	corss := cors.New(cors.Config{
		AllowOrigins:     cf.ApiConfig.CorsAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 允许携带凭证
		MaxAge:           12 * time.Hour, // 预检请求缓存时间
	})
	return corss
}
