package blog

import (
	"blog-go/database"
	"blog-go/log"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EnderCHX/chx-tools-go/encrypt"
	"go.uber.org/zap"
)

type Passage struct {
	PassageId      string    `json:"passage_id" gorm:"primaryKey;type:char(32);not null"`
	AuthorUsername string    `json:"author_username" gorm:"index;not null"`
	Content        string    `json:"content"`
	Title          string    `json:"title" gorm:"index;not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
	Deleted        bool      `json:"deleted"`
}

type Option func(*Passage)

func NewPassage(option ...Option) *Passage {
	passage := &Passage{}
	for _, opt := range option {
		opt(passage)
	}
	return passage
}

func WithPassageId(passageid string) Option {
	return func(p *Passage) {
		p.PassageId = passageid
	}
}

func WithTitle(title string) Option {
	return func(p *Passage) {
		p.Title = title
	}
}

func WithAuthorUsername(authorUsername string) Option {
	return func(p *Passage) {
		p.AuthorUsername = authorUsername
	}
}

func WithContent(content string) Option {
	return func(p *Passage) {
		p.Content = content
	}
}

func WithCreatedAt(createdAt time.Time) Option {
	return func(p *Passage) {
		p.CreatedAt = createdAt
	}
}

func WithUpdatedAt(updatedAt time.Time) Option {
	return func(p *Passage) {
		p.UpdatedAt = updatedAt
	}
}

func WithDeletedAt(deletedAt time.Time) Option {
	return func(p *Passage) {
		p.DeletedAt = deletedAt
	}
}

func WithDeleted(deleted bool) Option {
	return func(p *Passage) {
		p.Deleted = deleted
	}
}
func GetPassages() ([]Passage, error) {
	db := database.GetDB()
	var passages []Passage
	err := db.Select("passage_id", "title", "author_username", "created_at", "updated_at").Where("deleted = ?", 0).Find(&passages).Error
	if err != nil {
		return nil, err
	}
	return passages, nil
}

func GetPassagesCache() ([]Passage, error) {
	rdb, rctx := database.GetRedisClient()

	passagelist, err := rdb.Get(rctx, "blog:posts").Result()
	if err != nil {
		return nil, err
	}

	var passages []Passage
	err = json.Unmarshal([]byte(passagelist), &passages)
	if err != nil {
		return nil, err
	}
	return passages, nil
}

func UpdatePassagesCache(passages []Passage) error {
	valueToJson, err := json.Marshal(passages)
	if err != nil {
		return err
	}
	rdb, rctx := database.GetRedisClient()
	err = rdb.Set(rctx, "blog:posts", valueToJson, time.Hour*1).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetPassagesByTag(tagname string) ([]Passage, error) {
	db := database.GetDB()
	var passages []Passage
	result := db.Select("passages.passage_id", "passages.title", "passages.author_username", "passages.created_at", "passages.updated_at").
		Joins("JOIN passage_tags ON passage_tags.passage_id = passages.passage_id").
		Joins("JOIN tags ON tags.tag_id = passage_tags.tag_id").
		Where("tags.tag_name = ?", tagname).
		Where("passages.deleted = ?", 0).
		Find(&passages)

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("passage not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return passages, nil
}

func GetPassagesByIdCache(tagname string) ([]Passage, error) {
	rdb, rctx := database.GetRedisClient()
	passagelist, err := rdb.Get(rctx, "blog:tags:"+tagname).Result()
	if err != nil {
		return nil, err
	}
	var passages []Passage
	err = json.Unmarshal([]byte(passagelist), &passages)
	if err != nil {
		return nil, err
	}
	return passages, nil
}

func UpdatePassagesByIdCache(passages []Passage, tagName string) error {
	valueToJson, err := json.Marshal(passages)
	if err != nil {
		return err
	}
	rdb, rctx := database.GetRedisClient()
	err = rdb.Set(rctx, "blog:tags:"+tagName, valueToJson, time.Hour*1).Err()
	if err != nil {
		return err
	}
	return nil
}

func (p *Passage) CreateTable() {
	db := database.GetDB()
	err := db.AutoMigrate(&Passage{})
	if err != nil {
		log.Logger.Error("[MySQL]创建 passage 表失败", zap.String("error", err.Error()))
	} else {
		log.Logger.Info("[MySQL]创建 passage 表成功")
	}
}
func (p *Passage) GenPassageId() *Passage {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	p.PassageId = encrypt.Md5(p.AuthorUsername + p.CreatedAt.String() + p.Title)
	return p
}

func (p *Passage) Create() error {
	db := database.GetDB()
	err := db.Create(p).Error
	if err != nil {
		log.Logger.Error("[MySQL]创建 passage 失败", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (p *Passage) Update() error {
	db := database.GetDB()
	err := db.Model(p).Where("passage_id = ?", p.PassageId).Updates(p).Error
	if err != nil {
		log.Logger.Error("[MySQL]更新 passage 失败", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (p *Passage) Delete() error {
	db := database.GetDB()
	err := db.Model(p).Where("passage_id = ?", p.PassageId).Delete(p).Error
	if err != nil {
		log.Logger.Error("[MySQL]删除 passage 失败", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (p *Passage) Get() error {
	db := database.GetDB()
	err := db.Where("passage_id = ?", p.PassageId).Where("deleted = ?", 0).First(p).Error
	if err != nil {
		log.Logger.Error("[MySQL]获取 passage 失败", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (p *Passage) GetCache() error {
	rdb, rctx := database.GetRedisClient()
	passage, err := rdb.Get(rctx, "blog:posts:"+p.PassageId).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(passage), p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Passage) UpdateCache() error {
	valueToJson, err := json.Marshal(p)
	if err != nil {
		return err
	}
	rdb, rctx := database.GetRedisClient()
	err = rdb.Set(rctx, "blog:posts:"+p.PassageId, valueToJson, time.Hour*1).Err()
	if err != nil {
		return err
	}
	return nil
}

type PassageReqBody struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
