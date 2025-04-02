package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type VisitRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (v *VisitRepositoryImpl) AutoMigrate() {
	v.db.AutoMigrate(&entity.Visit{})
}

func (v *VisitRepositoryImpl) SetDb(db *gorm.DB) {
	v.db = db
}

func (v *VisitRepositoryImpl) SetRedis(rdb *cache.Redis) {
	v.rdb = rdb
}

func (v *VisitRepositoryImpl) SaveRecord(visit entity.Visit) error {
	return v.db.Create(&visit).Error
}
