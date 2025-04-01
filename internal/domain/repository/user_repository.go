package repository

import "blog-go/internal/domain/entity"

type UserRepository interface {
	SaveEmail(user entity.User) error
	GetEmail(username string) (string, error)
}
