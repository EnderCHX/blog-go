package repository

import (
	"blog-go/infrastructure/cache"

	"gorm.io/gorm"
)

type PassageTagsRepository interface {
	InitDb(db *gorm.DB, redis *cache.Redis)
	GetPassageTags(passageId string) ([]string, error)
	GetTagPassages(tagName string) ([]string, error)
}
