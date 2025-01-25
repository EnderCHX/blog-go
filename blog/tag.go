package blog

import (
	"blog-go/database"
	"blog-go/log"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type Tag struct {
	TagId   int    `json:"tag_id" gorm:"primaryKey"`
	TagName string `json:"tag_name" gorm:"index;unique;not null"`
}

func (t *Tag) CreateTable() {
	db := database.GetDB()
	err := db.AutoMigrate(&Tag{})
	if err != nil {
		log.Logger.Error("[MySQL]创建 tag 表失败", zap.String("error", err.Error()))
	} else {
		log.Logger.Info("[MySQL]创建 tag 表成功")
	}
}

func GetTags() ([]Tag, error) {
	db := database.GetDB()
	var tags []Tag
	err := db.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, err
}

func GetTagsCache() ([]Tag, error) {
	rdb, rctx := database.GetRedisClient()

	result, err := rdb.HGetAll(rctx, "blog:tags").Result()
	if len(result) == 0 {
		return nil, fmt.Errorf("no tags in cache")
	}

	if err != nil {
		return nil, err
	}

	var tags []Tag

	for k, v := range result {
		var tag Tag
		tag.TagName = k
		v_, _ := strconv.Atoi(v)
		tag.TagId = v_
		tags = append(tags, tag)
	}

	return tags, err
}

func UpdateTagsCache(tags []Tag) error {
	rdb, rctx := database.GetRedisClient()

	for _, tag := range tags {
		err := rdb.HSet(rctx, "blog:tags", tag.TagName, tag.TagId).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
