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
	posts.GET("/", controller.GetPassages)                   //获取所有文章
	posts.GET("/:tagname", controller.GetPassagesByTag)      //获取标签下的文章
	posts.GET("/tags", controller.GetTags)                   //获取所有标签
	posts.GET("/tags/:tagname", controller.GetPassagesByTag) //获取标签下的文章

	//获取文章
	post := r.Group("/post")
	post.GET("/:id", controller.GetPassageById)      //通过ID获取文章
	post.GET("/:id/tags", controller.GetPassageTags) //获取文章标签

	tag := r.Group("/tags")
	tag.GET("/", controller.GetTags)                  //获取所有标签
	tag.GET("/:tagname", controller.GetPassagesByTag) //获取标签下的文章

	r.Run(config_.Host + ":" + config_.Port)
}
