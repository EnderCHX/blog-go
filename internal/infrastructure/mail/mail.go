package mail

import "gopkg.in/gomail.v2"

type Mail struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewMail(host, port, username, password string) *Mail {
	return &Mail{Host: host, Port: port, Username: username, Password: password}
}

func (m *Mail) SendMail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Username)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	d := gomail.NewDialer(m.Host, 587, m.Username, m.Password)
	return d.DialAndSend(msg)
}
