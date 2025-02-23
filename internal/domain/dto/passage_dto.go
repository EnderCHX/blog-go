package dto

import "time"

type PassageDTO struct {
	PassageId      string    `json:"passage_id"`
	AuthorUsername string    `json:"author_username"`
	Content        string    `json:"content"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
	Deleted        bool      `json:"deleted"`
}

type NewPassageDTO struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
