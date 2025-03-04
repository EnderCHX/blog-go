package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (u *UserRepositoryImpl) AutoMigrate() {
	u.db.AutoMigrate(&entity.User{})
}

func (u *UserRepositoryImpl) SetDb(db *gorm.DB) {
	u.db = db
}

func (u *UserRepositoryImpl) SetRedis(rdb *cache.Redis) {
	u.rdb = rdb
}

func (u *UserRepositoryImpl) SaveEmail(user entity.User) error {
	u.rdb.Set("blog:user:mail:"+user.Username, user.Email, 0)
	return u.db.Create(&user).Error
}

func (u *UserRepositoryImpl) GetEmail(username string) (string, error) {
	email, err := u.rdb.Get("blog:user:mail:" + username)
	if err != nil {
		err2 := u.db.Where("username = ?", username).Pluck("email", &email).Error
		if err2 != nil {
			return "", err2
		}

		if err == redis.Nil {
			go func() {
				u.rdb.Set("blog:user:mail:"+username, email, 0)
			}()
		}
	}
	return email, nil
}
