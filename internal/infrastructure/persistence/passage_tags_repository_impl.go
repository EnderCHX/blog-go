package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"

	"gorm.io/gorm"
)

type PassageTagsRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (p *PassageTagsRepositoryImpl) InitDb(db *gorm.DB, redis *cache.Redis) {
	p.db = db
	p.rdb = redis
	p.db.AutoMigrate(&entity.PassageTags{})
}

func (p *PassageTagsRepositoryImpl) GetPassageTags(passageId string) ([]string, error) {
	var tags []string
	err := p.db.Model(&entity.PassageTags{}).
		Select("tag_name").
		Where("passage_id = ?", passageId).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
func (p *PassageTagsRepositoryImpl) GetTagPassages(tagName string) ([]string, error) {
	var passagesId []string
	err := p.db.Model(&entity.PassageTags{}).
		Select("passage_id").
		Where("tag_name = ?", tagName).
		Find(&passagesId).Error
	if err != nil {
		return nil, err
	}
	return passagesId, err
}
