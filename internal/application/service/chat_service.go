package service

import "blog-go/internal/domain/entity"

type ChatService interface {
	AddChat(chat *entity.Chat) error
	GetAll() ([]entity.Chat, error)
	DeleteChat(chatId string) error
}
