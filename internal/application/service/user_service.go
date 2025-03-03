package service

import "blog-go/internal/domain/entity"

type UserService interface {
	GetUserInfo(username string) (entity.User, error)
}
