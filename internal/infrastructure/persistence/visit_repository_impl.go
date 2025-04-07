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

func (v *VisitRepositoryImpl) CountAll() (int64, error) {
	var count int64
	return count, v.db.Model(&entity.Visit{}).Count(&count).Error
}

func (v *VisitRepositoryImpl) CountPath(path string) (int64, error) {
	var count int64
	return count, v.db.Model(&entity.Visit{}).Where("path = ?", path).Count(&count).Error
}

func (v *VisitRepositoryImpl) CountIp(ip string) (int64, error) {
	var count int64
	return count, v.db.Model(&entity.Visit{}).Where("ip = ?", ip).Count(&count).Error
}

func (v *VisitRepositoryImpl) CountBrowser(browser string) (int64, error) {
	var count int64
	return count, v.db.Model(&entity.Visit{}).Where("browser = ?", browser).Count(&count).Error
}

func (v *VisitRepositoryImpl) CountOs(os string) (int64, error) {
	var count int64
	return count, v.db.Model(&entity.Visit{}).Where("os = ?", os).Count(&count).Error
}

func (v *VisitRepositoryImpl) CountPlatform(platform string) (int64, error) {
	var count int64
	return count, v.db.Model(&entity.Visit{}).Where("platform = ?", platform).Count(&count).Error
}
