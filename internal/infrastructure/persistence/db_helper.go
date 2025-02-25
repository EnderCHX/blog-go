package persistence

import (
	"blog-go/internal/infrastructure/cache"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbHepler struct {
	PassageRepsitory     PassageRepositoryImpl
	PassageTagsRepsitory PassageTagsRepositoryImpl
	TagsRepository       TagsRepositoryImpl
	db                   *gorm.DB
	redis                *cache.Redis
}

func (h *DbHepler) InitDbHepler(mysql *gorm.DB, redis *cache.Redis) error {
	h.db = mysql
	h.redis = redis

	v := reflect.ValueOf(h).Elem()
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		if !field.CanInterface() || fieldType.Name == "db" || fieldType.Name == "redis" {
			continue
		}

		val := field

		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			continue
		}

		if val.Addr().MethodByName("SetDb").IsValid() {
			val.Addr().MethodByName("SetDb").Call([]reflect.Value{reflect.ValueOf(mysql)})
		}

		if val.Addr().MethodByName("SetRedis").IsValid() {
			val.Addr().MethodByName("SetRedis").Call([]reflect.Value{reflect.ValueOf(redis)})
		}
	}

	return nil
}

func (h *DbHepler) AutoMigrate() {
	v := reflect.ValueOf(h).Elem() // 获取结构体反射对象
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		// 跳过非公开字段和特定字段（如db/redis）
		if !field.CanInterface() || fieldType.Name == "db" || fieldType.Name == "redis" {
			continue
		}

		// 处理指针类型：获取实际指向的值
		val := field
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// 确认是否为结构体类型
		if val.Kind() != reflect.Struct {
			continue
		}

		// 查找并调用AutoMigrate方法
		method := val.Addr().MethodByName("AutoMigrate")
		if method.IsValid() && method.Type().NumIn() == 0 {
			method.Call(nil)
		}
	}
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
