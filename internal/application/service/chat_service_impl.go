package service

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"
)

type ChatServiceImpl struct {
	db *persistence.DbHelper
}

func NewChatService(db *persistence.DbHelper) ChatService {
	return &ChatServiceImpl{db: db}
}
func (c *ChatServiceImpl) AddChat(chat *entity.Chat) error {
	return c.db.ChatRepository.AddChat(chat)
}

func (c *ChatServiceImpl) GetAll() ([]entity.Chat, error) {
	return c.db.ChatRepository.GetAll()
}

func (c *ChatServiceImpl) DeleteChat(chatId string) error {
	return c.db.ChatRepository.DeleteChat(chatId)
}
