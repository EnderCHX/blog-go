package main

import (
	"blog-go/api"
	"blog-go/blog"
	"blog-go/config"
	"blog-go/database"
	"blog-go/log"
)

func main() {
	config.Setup()
	log.Setup()
	database.Setup()
	blog.Setup()
	api.StartApi()
}
