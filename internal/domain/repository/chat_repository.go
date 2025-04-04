package repository

import "blog-go/internal/domain/entity"

type ChatRepository interface {
	AddChat(chat *entity.Chat) error
	GetAll() ([]entity.Chat, error)
	GetChatsByUser(username string) ([]entity.Chat, error)
	GetChatById(chatId string) (*entity.Chat, error)
	DeleteChat(chatId string) error
}
