package repository

import (
	"blog-go/internal/domain/entity"
)

type TagsRepository interface {
	GetTags() ([]entity.Tag, error)
	GetTagName(id int) (string, error)
	GetTagId(tagname string) (int, error)
	AddTags(tags []entity.Tag) error
	DeleteTags(tags []entity.Tag) error
}
