package blog

import (
	"time"

	"github.com/EnderCHX/chx-tools-go/encrypt"
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

func (p *Passage) GetPassageId() *Passage {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	p.PassageId = encrypt.Md5(p.AuthorUsername + p.CreatedAt.String() + p.Title)
	return p
}
