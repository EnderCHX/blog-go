package blog

import (
	"blog-go/database"
	"blog-go/log"

	"go.uber.org/zap"
)

type PassageTags struct {
	PassageId string `json:"passage_id" gorm:"primaryKey"`
	TagId     int    `json:"tag_id" gorm:"not null"`
}

func (pt *PassageTags) CreateTable() {
	db := database.GetDB()
	err := db.AutoMigrate(&PassageTags{})
	if err != nil {
		log.Logger.Error("[MySQL]创建 passage_tags 表失败", zap.String("error", err.Error()))
	} else {
		log.Logger.Info("[MySQL]创建 passage_tags 表成功")
	}
}
