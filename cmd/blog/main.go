package main

import (
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/config"
	"blog-go/internal/infrastructure/log"
	"blog-go/internal/infrastructure/persistence"
	"fmt"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

func main() {
	config.Setup("./config.json")
	cf := config.GetConfig()

	log.Setup(cf.LogCongfig.LogPath, cf.LogCongfig.LogLevel)
	logger := log.GetLogger()

	dbhelper := persistence.DbHepler{}
	mysql, err := persistence.NewMySQL(
		cf.MySQLConfig.Host,
		cf.MySQLConfig.Port,
		cf.MySQLConfig.Username,
		cf.MySQLConfig.Password,
		cf.MySQLConfig.DBName,
		&log.GormLogger{Logger: logger, LogLevel: gormlogger.Info})
	if err != nil {
		logger.Error("数据库连接失败")
	}
	redis := cache.NewRedis(cf.RedisConfig.Host, cf.RedisConfig.Port, cf.RedisConfig.Username, cf.RedisConfig.Password, cf.RedisConfig.DB)
	dbhelper.InitDbHepler(mysql, redis)
	dbhelper.AutoMigrate()
	passages, _ := dbhelper.PassageRepository.GetPassages()
	logger.Info(fmt.Sprintf("%v", passages))

	tags, _ := dbhelper.TagsRepository.GetTags()
	logger.Info(fmt.Sprintf("%v", tags))

	passages2, _ := dbhelper.PassageTagsRepository.GetTagPassages("22")
	logger.Info(fmt.Sprintf("%v", passages2))

	tags2, _ := dbhelper.PassageTagsRepository.GetPassageTags("7de19be4d5898f4abc022206954a8ad6")
	logger.Info(fmt.Sprintf("%v", tags2))

	id, _ := dbhelper.TagsRepository.GetTagId("openwrt")
	tag, _ := dbhelper.TagsRepository.GetTagName(id)
	logger.Info(fmt.Sprintf("%v %v", id, tag))

	start, _ := time.Parse("2006-01-02", "2025-01-01")
	passages3, err := dbhelper.PassageRepository.GetPassagesByDate(start, time.Now())
	logger.Info(fmt.Sprintf("%v", passages3))
	if err != nil {
		logger.Error(err.Error())
	}
	for {
	}
}
