package entity

import "time"

type Comment struct {
	PassageId string    `json:"passage_id" gorm:"index;not null"`
	Username  string    `json:"username" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Deleted   bool      `json:"deleted"`
}
