package service

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"
	"time"
)

type PassageServiceImpl struct {
	db persistence.DbHepler
}

func NewPassageServiceImpl(db persistence.DbHepler) PassageService {
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
