package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"

	"gorm.io/gorm"
)

type TagsRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (t *TagsRepositoryImpl) AutoMigrate() {
	t.db.AutoMigrate(&entity.Tag{})
}

func (p *TagsRepositoryImpl) SetDb(db *gorm.DB) {
	p.db = db
}

func (p *TagsRepositoryImpl) SetRedis(rdb *cache.Redis) {
	p.rdb = rdb
}
func (t *TagsRepositoryImpl) GetTags() ([]entity.Tag, error) {
	var tags []entity.Tag
	err := t.db.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, err
}
func (t *TagsRepositoryImpl) GetTagName(id int) (string, error) {
	var tagname string
	err := t.db.Model(&entity.Tag{}).Select("tag_name").Where("tag_id = ?", id).Find(&tagname).Error
	if err != nil {
		return "", err
	}
	return tagname, err
}
func (t *TagsRepositoryImpl) GetTagId(tagname string) (int, error) {
	var tagId int
	err := t.db.Model(&entity.Tag{}).Select("tag_id").Where("tag_name = ?", tagname).Find(&tagId).Error
	if err != nil {
		return -1, err
	}
	return tagId, nil
}
func (t *TagsRepositoryImpl) AddTags(tags []entity.Tag) error {
	return t.db.Create(&tags).Error
}
func (t *TagsRepositoryImpl) DeleteTags(tags []entity.Tag) error {
	return t.db.Delete(&tags).Error
}
