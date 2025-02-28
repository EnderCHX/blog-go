package service

import (
	"blog-go/internal/domain/entity"
	"time"
)

type PassageService interface {
	GetPassages() ([]entity.Passage, error)
	GetPassagesByTag(tagname string) ([]string, error)
	GetPassagesByDate(start, end time.Time) ([]string, error)
	GetPassageById(id string) (entity.Passage, error)

	AddPassage(passage entity.Passage) error
	AddPassageTags(passageId string, tags []string) error
	AddTags(tags []string) error
}
