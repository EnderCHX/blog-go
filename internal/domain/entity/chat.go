package entity

import (
	"github.com/EnderCHX/chx-tools-go/encrypt"
	"time"
)

type Chat struct {
	ChatId      string    `json:"chat_id" gorm:"primaryKey;type:char(32);not null"`
	ReplyChatId string    `json:"reply_chat_id"`
	Username    string    `json:"username" gorm:"index;not null"`
	Content     string    `json:"content" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	Deleted     bool      `json:"deleted"`
}

func (c *Chat) GenChatId() *Chat {
	c.ChatId = encrypt.Md5(c.Username + c.Content + c.CreatedAt.String())
	return c
}
