package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"

	"gorm.io/gorm"
)

type PassageRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (p *PassageRepositoryImpl) AutoMigrate() {
	p.db.AutoMigrate(&entity.Passage{})
}

func (p *PassageRepositoryImpl) SetDb(db *gorm.DB) {
	p.db = db
}

func (p *PassageRepositoryImpl) SetRedis(rdb *cache.Redis) {
	p.rdb = rdb
}
func (p *PassageRepositoryImpl) GetPassage(id string) (entity.Passage, error) {
	var passage entity.Passage
	err := p.db.Where("passage_id = ?", id).Where("deleted = ?", 0).First(&passage).Error
	if err != nil {
		return passage, err
	}
	return passage, nil
}
func (p *PassageRepositoryImpl) GetPassages() ([]entity.Passage, error) {
	var passages []entity.Passage
	err := p.db.
		Select("passage_id",
			"title",
			"author_username",
			"created_at",
			"updated_at").
		Where("deleted = ?", 0).
		Find(&passages).Error

	if err != nil {
		return nil, err
	}
	return passages, nil
}
func (p *PassageRepositoryImpl) AddPassage(passage entity.Passage) error {
	return p.db.Create(&passage).Error
}
func (p *PassageRepositoryImpl) UpdatePassage(passage entity.Passage) error {
	return p.db.Model(&passage).Where("passage_id = ?", passage.PassageId).Updates(&passage).Error
}
func (p *PassageRepositoryImpl) DeletePassage(id int) error {
	return p.db.Model(&entity.Passage{}).Where("passage_id = ?", id).Update("deleted", true).Error
}
