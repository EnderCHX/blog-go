package database

import (
	"blog-go/config"
	"blog-go/log"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func Setup() {
	log.Logger.Info("[MySQL]数据库连接中...")
	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("数据库连接失败", zap.Any("err", err))
		}
	}()

	gormLogger := &log.GormLogger{
		Logger:   log.Logger,
		LogLevel: logger.Info,
	}

	mysqlConfig := config.ConfigContext.MySQLConfig
	dsn := mysqlConfig.Username + ":" + mysqlConfig.Password + "@tcp(" + mysqlConfig.Host + ":" + mysqlConfig.Port + ")/" + mysqlConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic("failed to connect database")
	} else {
		log.Logger.Info("[MySQL]数据库连接成功")
	}
	DB = db
	GetRedisClient()
}

func GetDB() *gorm.DB {
	if DB == nil {
		Setup()
	}
	return DB
}
