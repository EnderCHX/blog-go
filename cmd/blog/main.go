package main

import (
	"blog-go/internal/application/service"
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/config"
	"blog-go/internal/infrastructure/log"
	"blog-go/internal/infrastructure/mail"
	"blog-go/internal/infrastructure/persistence"
	"blog-go/internal/interfaces"
	"blog-go/internal/interfaces/handle"
	"blog-go/internal/interfaces/midware"
	"blog-go/internal/interfaces/route"

	"github.com/gin-gonic/gin"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	config.Setup("./config.json")
	cf := config.GetConfig()

	log.Setup(cf.LogCongfig.LogPath, cf.LogCongfig.LogLevel)
	logger := log.GetLogger()

	dbHelper := &persistence.DbHelper{}
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
	dbHelper.InitDbHelper(mysql, redis)
	dbHelper.AutoMigrate()

	m := mail.NewMail(cf.MailConfig.Host, cf.MailConfig.Port, cf.MailConfig.Username, cf.MailConfig.Password)

	passageHandle := handle.NewPassageHandle(service.NewPassageServiceImpl(dbHelper), logger)
	tagHandle := handle.NewTagHandle(service.NewTagServiceImpl(dbHelper), logger)
	commentHandle := handle.NewCommentHandle(service.NewCommentServiceImpl(dbHelper), service.NewUserServiceImpl(dbHelper, m, cf.ApiConfig.UserApiUrl), logger)
	chatHandle := handle.NewChatHandle(service.NewChatService(dbHelper), logger)
	visitHandle := handle.NewVisitHandle(service.NewVisitServiceImpl(dbHelper))

	server := interfaces.NewHttpServer(cf.ApiConfig.Host, cf.ApiConfig.Port, cf.ApiConfig.Mode,
		[]gin.HandlerFunc{log.GinZapLogger(), gin.Recovery(), midware.CountVisitor(dbHelper, cf), midware.Cors(cf), midware.Auth(cf)},
		route.NewPassageRouteRegister(passageHandle),
		route.NewTagRouteRegister(tagHandle),
		route.NewCommentRouteRegister(commentHandle),
		route.NewChatRouteRegister(chatHandle),
		route.NewVisitRoute(visitHandle),
	)
	//go m.SendMail("c@chxc.cc", "服务器启动", "api服务器启动")
	server.Start()
}
