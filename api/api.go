package api

import (
	"blog-go/config"
	"blog-go/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	config_ := config.ConfigContext.ApiConfig
	gin.SetMode(config_.Mode)

	r := gin.New()

	r.Use(log.GinZapLogger(), gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to my API",
		})
	})

	r.Run(config_.Host + ":" + config_.Port)
}
