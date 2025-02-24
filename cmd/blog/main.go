package main

import (
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/config"
	"blog-go/internal/infrastructure/log"
	"blog-go/internal/infrastructure/persistence"
	"fmt"

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
	passages, _ := dbhelper.PassageRepsitory.GetPassages()
	for _, passage := range passages {
		fmt.Println(passage)
	}
	tags, _ := dbhelper.TagsRepository.GetTags()

	for _, tag := range tags {
		fmt.Println(tag)
	}

	passages2, _ := dbhelper.PassageTagsRepsitory.GetTagPassages(9)
	fmt.Println(passages2)

}
