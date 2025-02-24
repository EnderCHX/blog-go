package repository

import (
	"blog-go/internal/domain/entity"
)

type PassageRepository interface {
	GetPassage(id int) (entity.Passage, error)
	GetPassages() ([]entity.Passage, error)
	AddPassage(passage entity.Passage) error
	UpdatePassage(passage entity.Passage) error
	DeletePassage(id int) error
}
