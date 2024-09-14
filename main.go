package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	r := gin.Default()

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

	r.POST("/get", login)

	r.Run(":8080")
}

func login(c *gin.Context) {
	user := User{}
	c.BindJSON(&user)
	c.JSON(http.StatusOK, gin.H{
		"message":  "test",
		"username": user.Username,
		"password": user.Password,
	})
}
