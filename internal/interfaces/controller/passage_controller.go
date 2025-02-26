package controller

import "blog-go/internal/interfaces"

type PassageController struct {
	httpserver *interfaces.HttpServer
}

func NewPassageController(httpserver *interfaces.HttpServer) *PassageController {
	return &PassageController{
		httpserver: httpserver,
	}
}
