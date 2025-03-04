package service

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/mail"
	"blog-go/internal/infrastructure/persistence"
)

type UserServiceImpl struct {
	db         *persistence.DbHelper
	mail       *mail.Mail
	userApiUrl string
}

func NewUserServiceImpl(db *persistence.DbHelper, mail *mail.Mail, userApiUrl string) UserService {
	return &UserServiceImpl{
		db:         db,
		mail:       mail,
		userApiUrl: userApiUrl,
	}
}

func (u *UserServiceImpl) ServiceName() string {
	return "UserService"
}

func (u *UserServiceImpl) SavaEmail(username, accessToken string) error {
	mail, _ := u.db.UserRepository.GetEmail(username)
	if mail != "" {
		return nil
	}
	user := entity.NewUser(username, "")
	var err error
	user.Email, err = user.GetEmail(u.userApiUrl, accessToken)
	if err != nil {
		return err
	}
	return u.db.UserRepository.SaveEmail(*user)
}

func (u *UserServiceImpl) SendEmail(username, subject, body string) error {
	email, err := u.db.UserRepository.GetEmail(username)
	if err != nil {
		return err
	}
	return u.mail.SendMail(email, subject, body)
}
