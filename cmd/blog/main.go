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

	dbhepler := persistence.DbHepler{}
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
	dbhepler.InitDbHepler(mysql, redis)
	passages, _ := dbhepler.PassageRepsitory.GetPassages()
	for _, passage := range passages {
		fmt.Println(passage)
	}
}
