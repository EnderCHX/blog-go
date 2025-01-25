package blog

import (
	"blog-go/database"
	"blog-go/log"

	"go.uber.org/zap"
)

type UserPermissions struct {
	Username    string `json:"username" gorm:"primaryKey"`
	Permissions string `json:"permissions"`
}

func (up *UserPermissions) CreateTable() {
	db := database.GetDB()
	err := db.AutoMigrate(&UserPermissions{})
	if err != nil {
		log.Logger.Error("[MySQL]创建 user_permissions 表失败", zap.String("error", err.Error()))
		return
	}
	log.Logger.Info("[MySQL]创建 user_permissions 表成功")
}
