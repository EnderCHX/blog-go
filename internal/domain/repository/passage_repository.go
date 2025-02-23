package repository

import (
	"blog-go/domain/entity"
	"blog-go/infrastructure/cache"

	"gorm.io/gorm"
)

type PassageRepository interface {
	InitDb(db *gorm.DB, redis *cache.Redis)
	GetPassage(id int) (entity.Passage, error)
	GetPassages() ([]entity.Passage, error)
	AddPassage(passage entity.Passage) error
	UpdatePassage(passage entity.Passage) error
	DeletePassage(id int) error
}
