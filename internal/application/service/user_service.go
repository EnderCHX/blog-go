package service

type UserService interface {
	SavaEmail(username, accessToken string) error
	SendEmail(username, subject, body string) error
}
