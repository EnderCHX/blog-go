package api

import (
	"blog-go/api/controller"
	"blog-go/api/midware"
	"blog-go/config"
	"blog-go/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartApi() {
	config_ := config.ConfigContext.ApiConfig
	gin.SetMode(config_.Mode)

	r := gin.New()

	r.Use(log.GinZapLogger(), gin.Recovery(), midware.Auth())

	r.StaticFile("/favicon.ico", "./favicon.ico")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ciallo World",
		})
	})

	//新建文章
	r.POST("/newpost", controller.NewPost)

	//文章列表
	posts := r.Group("/posts")
	{
		posts.GET("/", controller.GetPassages)
		posts.GET("/:tagname", controller.GetPassagesByTag)
	}

	//获取文章
	post := r.Group("/post")
	{
		post.GET("/:id", controller.GetPassageById)
	}

	tag := r.Group("/tag")
	{
		//获取标签
		tag.GET("/", controller.GetTags)
		tag.GET("/:tagname", controller.GetPassagesByTag)
	}

	r.Run(config_.Host + ":" + config_.Port)
}
