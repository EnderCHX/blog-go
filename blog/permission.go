package blog

import (
	"blog-go/database"
	"blog-go/log"

	"go.uber.org/zap"
)

type Permission struct {
	PermissionId int    `json:"permission_id" gorm:"primaryKey"`
	Permission   string `json:"permission"`
}

var Permissions = []Permission{
	{PermissionId: 0, Permission: "access"},  //访问
	{PermissionId: 1, Permission: "read"},    //读文章
	{PermissionId: 2, Permission: "edit"},    //写	文章
	{PermissionId: 3, Permission: "comment"}, //评论
}

func (p *Permission) CreateTable() {
	db := database.GetDB()
	err := db.AutoMigrate(&Permission{})
	if err != nil {
		log.Logger.Error("[MySQL]创建 permission 表失败", zap.String("error", err.Error()))
		return
	}
	log.Logger.Info("[MySQL]创建 permission 表成功")

	db.Create(&Permissions)
}
