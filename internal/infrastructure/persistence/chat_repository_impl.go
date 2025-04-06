package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type ChatRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (c *ChatRepositoryImpl) AutoMigrate() {
	c.db.AutoMigrate(&entity.Chat{})
}

func (c *ChatRepositoryImpl) SetDb(db *gorm.DB) {
	c.db = db
}

func (c *ChatRepositoryImpl) SetRedis(rdb *cache.Redis) {
	c.rdb = rdb
}

func (c *ChatRepositoryImpl) AddChat(chat *entity.Chat) error {
	return c.db.Create(chat).Error
}

func (c *ChatRepositoryImpl) GetAll() ([]entity.Chat, error) {
	var chats []entity.Chat
	err := c.db.Limit(20).Where("deleted = ?", false).Order("created_at desc").Find(&chats).Error

	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (c *ChatRepositoryImpl) GetChatsByUser(username string) ([]entity.Chat, error) {
	var chats []entity.Chat
	err := c.db.Where("username = ?", username).Find(&chats).Error

	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (c *ChatRepositoryImpl) GetChatById(chatId string) (*entity.Chat, error) {
	var chat *entity.Chat
	err := c.db.Where("chat_id = ?", chatId).Find(&chat).Error

	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *ChatRepositoryImpl) DeleteChat(chatId string) error {
	return c.db.Where("chat_id = ?", chatId).Update("deleted = ?", true).Error
}
