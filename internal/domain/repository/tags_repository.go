package repository

import (
	"blog-go/domain/entity"
	"blog-go/infrastructure/cache"

	"gorm.io/gorm"
)

type TagsRepository interface {
	InitDb(db *gorm.DB, redis *cache.Redis)
	GetTags() ([]entity.Tag, error)
	GetTagName(id int) (string, error)
	GetTagId(tagname string) (int, error)
	AddTags(tags []entity.Tag) error
	DeleteTags(tags []entity.Tag) error
}
