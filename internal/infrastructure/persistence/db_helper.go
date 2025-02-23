package persistence

import (
	"blog-go/internal/infrastructure/cache"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbHepler struct {
	PassageRepsitory     PassageRepositoryImpl
	PassageTagsRepsitory PassageTagsRepositoryImpl
	TagsRepository       PassageTagsRepositoryImpl
	db                   *gorm.DB
	redis                *cache.Redis
}

func (d *DbHepler) InitDbHepler(mysql *gorm.DB, redis *cache.Redis) error {
	d.db = mysql
	d.redis = redis
	{
		d.PassageRepsitory.InitDb(mysql, redis)
		d.PassageTagsRepsitory.InitDb(mysql, redis)
		d.TagsRepository.InitDb(mysql, redis)
	}
	return nil
}

func NewMySQL(host, port, username, password, db_name string, logger logger.Interface) (*gorm.DB, error) {
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
