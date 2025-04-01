package service

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"
)

type TagServiceImpl struct {
	db *persistence.DbHelper
}

func (t *TagServiceImpl) ServiceName() string {
	return "TagService"
}

func NewTagServiceImpl(db *persistence.DbHelper) TagService {
	return &TagServiceImpl{
		db: db,
	}
}

func (t *TagServiceImpl) GetTags() ([]string, error) {
	tags, err := t.db.TagsRepository.GetTags()
	if err != nil {
		return nil, err
	}

	var tagstings []string
	for _, tag := range tags {
		tagstings = append(tagstings, tag.TagName)
	}

	return tagstings, nil
}

func (t *TagServiceImpl) GetPassageTags(passageId string) ([]string, error) {
	tags, err := t.db.PassageTagsRepository.GetPassageTags(passageId)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *TagServiceImpl) GetTagName(id int) (string, error) {
	name, err := t.db.TagsRepository.GetTagName(id)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (t *TagServiceImpl) GetTagId(tagname string) (int, error) {
	id, err := t.db.TagsRepository.GetTagId(tagname)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (t *TagServiceImpl) AddTags(tags []string) error {
	var tagse []entity.Tag
	for _, tag := range tags {
		tagse = append(tagse, entity.Tag{TagName: tag})
	}
	return t.db.TagsRepository.AddTags(tagse)
}

func (t *TagServiceImpl) DeleteTags(tags []string) error {
	var tagse []entity.Tag
	for _, tag := range tags {
		tagse = append(tagse, entity.Tag{TagName: tag})
	}
	return t.db.TagsRepository.DeleteTags(tagse)
}
