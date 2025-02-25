package service

import (
	"blog-go/internal/domain/entity"
	"time"
)

type PassageService interface {
	GetPassages() ([]entity.Passage, error)
	GetPassagesByTag(tagname string) ([]string, error)
	GetPassagesByDate(start, end time.Time) ([]string, error)
	GetPassageById() (entity.Passage, error)
}
