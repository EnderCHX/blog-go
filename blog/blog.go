package blog

import (
	"blog-go/log"
)

type BlogTable interface {
	CreateTable()
}

func Setup() {
	BlogTables := []BlogTable{
		&Passage{},
		&Tag{},
		&PassageTags{},
		&Permission{},
		&UserPermissions{},
	}
	for _, table := range BlogTables {
		table.CreateTable()
	}
}

func Hello() {
	log.Logger.Info("Hello")
}
