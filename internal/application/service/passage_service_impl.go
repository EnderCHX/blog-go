package service

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"
	"time"
)

type PassageServiceImpl struct {
	db *persistence.DbHepler
}

func NewPassageServiceImpl(db *persistence.DbHepler) PassageService {
	return &PassageServiceImpl{
		db: db,
	}
}

func (p *PassageServiceImpl) GetPassages() ([]entity.Passage, error) {
	return p.db.PassageRepository.GetPassages()
}

func (p *PassageServiceImpl) GetPassagesByTag(tagname string) ([]string, error) {
	passageIds, err := p.db.PassageTagsRepository.GetTagPassages(tagname)
	if err != nil {
		return nil, err
	}

	return passageIds, nil
}

func (p *PassageServiceImpl) GetPassagesByDate(start, end time.Time) ([]string, error) {
	return p.db.PassageRepository.GetPassagesByDate(start, end)
}

func (p *PassageServiceImpl) GetPassageById(id string) (entity.Passage, error) {
	return p.db.PassageRepository.GetPassage(id)
}

func (p *PassageServiceImpl) AddPassage(passage entity.Passage) error {
	passage.Deleted = false
	passage.GenPassageId()
	return p.db.PassageRepository.AddPassage(passage)
}

func (t *PassageServiceImpl) AddPassageTags(passageId string, tags []string) error {
	return t.db.PassageTagsRepository.AddPassageTags(passageId, tags)
}

func (t *PassageServiceImpl) AddTags(tags []string) error {
	var tagse []entity.Tag
	for _, tag := range tags {
		tagse = append(tagse, entity.Tag{TagName: tag})
	}
	return t.db.TagsRepository.AddTags(tagse)
}
