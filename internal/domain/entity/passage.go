package entity

import (
	"time"
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

type PassageOptions func(*Passage)

func NewPassage(option ...PassageOptions) *Passage {
	passage := &Passage{}
	for _, opt := range option {
		opt(passage)
	}
	return passage
}

func WithPassageId(passageid string) PassageOptions {
	return func(p *Passage) {
		p.PassageId = passageid
	}
}

func WithTitle(title string) PassageOptions {
	return func(p *Passage) {
		p.Title = title
	}
}

func WithAuthorUsername(authorUsername string) PassageOptions {
	return func(p *Passage) {
		p.AuthorUsername = authorUsername
	}
}

func WithContent(content string) PassageOptions {
	return func(p *Passage) {
		p.Content = content
	}
}

func WithCreatedAt(createdAt time.Time) PassageOptions {
	return func(p *Passage) {
		p.CreatedAt = createdAt
	}
}

func WithUpdatedAt(updatedAt time.Time) PassageOptions {
	return func(p *Passage) {
		p.UpdatedAt = updatedAt
	}
}

func WithDeletedAt(deletedAt time.Time) PassageOptions {
	return func(p *Passage) {
		p.DeletedAt = deletedAt
	}
}

func WithDeleted(deleted bool) PassageOptions {
	return func(p *Passage) {
		p.Deleted = deleted
	}
}
