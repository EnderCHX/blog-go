package service

import (
	"blog-go/internal/domain/entity"
	"time"
)

type PassageService interface {
	GetPassages() ([]entity.Passage, error)
	GetPassagesByTag(tagname string) ([]entity.Passage, error)
	GetPassagesByDate(start, end time.Time) ([]entity.Passage, error)
	GetPassageById() (entity.Passage, error)
}
