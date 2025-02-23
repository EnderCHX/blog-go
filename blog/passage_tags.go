package blog

import (
	"blog-go/database"
	"blog-go/log"
	"fmt"

	"go.uber.org/zap"
)

type PassageTags struct {
	PassageId string `json:"passage_id" gorm:"not null"`
	TagId     int    `json:"tag_id" gorm:"not null"`
}

func CreatePassageTags(passageTags []PassageTags) error {
	db := database.GetDB()
	err := db.Create(&passageTags).Error
	if err != nil {
		log.Logger.Error("[MySQL]创建 passage_tags 失败", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func GetPassageTags(passageId string) ([]string, error) {
	db := database.GetDB()
	var tags []string

	err := db.Table("passage_tags").
		Select("tag_name").
		Joins("JOIN tags ON tags.tag_id = passage_tags.tag_id").
		Where("passage_tags.passage_id = ?", passageId).
		Pluck("tags.tag_name", &tags).Error
	return tags, err
}

func GetPassageTagsCache(passageId string) ([]string, error) {
	rdb, rctx := database.GetRedisClient()
	tags, err := rdb.SMembers(rctx, "blog:posts:"+passageId+":tags").Result()
	if len(tags) == 0 {
		return nil, fmt.Errorf("no tags in cache")
	}
	return tags, err
}

func UpdatePassageTagsCache(passageId string, tags []string) error {
	rdb, rctx := database.GetRedisClient()
	err := rdb.SAdd(rctx, "blog:posts:"+passageId+":tags", tags).Err()
	return err
}
