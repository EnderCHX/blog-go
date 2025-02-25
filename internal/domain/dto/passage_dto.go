package dto

import (
	"blog-go/internal/domain/entity"
	"encoding/json"
	"time"
)

type PassageDTO struct {
	PassageId      string    `json:"passage_id"`
	AuthorUsername string    `json:"author_username"`
	Content        string    `json:"content"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PassageDTOOptions func(*PassageDTO)

func NewPassageDTO(option ...PassageDTOOptions) *PassageDTO {
	p := &PassageDTO{}
	for _, o := range option {
		o(p)
	}
	return p
}

func WithPassageId(passageId string) PassageDTOOptions {
	return func(p *PassageDTO) {
		p.PassageId = passageId
	}
}

func WithAuthorUsername(authorUsername string) PassageDTOOptions {
	return func(p *PassageDTO) {
		p.AuthorUsername = authorUsername
	}
}

func WithContent(content string) PassageDTOOptions {
	return func(p *PassageDTO) {
		p.Content = content
	}
}

func WithTitle(title string) PassageDTOOptions {
	return func(p *PassageDTO) {
		p.Title = title
	}
}

func WithCreatedAt(createdAt time.Time) PassageDTOOptions {
	return func(p *PassageDTO) {
		p.CreatedAt = createdAt
	}
}

func WithUpdatedAt(updatedAt time.Time) PassageDTOOptions {
	return func(p *PassageDTO) {
		p.UpdatedAt = updatedAt
	}
}

func (p *PassageDTO) ToJson() string {
	j, _ := json.Marshal(p)
	return string(j)
}

func (p *PassageDTO) ToPassage() entity.Passage {
	return *entity.NewPassage(
		entity.WithPassageId(p.PassageId),
		entity.WithAuthorUsername(p.AuthorUsername),
		entity.WithContent(p.Content),
		entity.WithTitle(p.Title),
		entity.WithCreatedAt(p.CreatedAt),
		entity.WithUpdatedAt(p.UpdatedAt),
	)
}

func PassageToDTO(passage entity.Passage) *PassageDTO {
	return NewPassageDTO(
		WithPassageId(passage.PassageId),
		WithAuthorUsername(passage.AuthorUsername),
		WithContent(passage.Content),
		WithTitle(passage.Title),
		WithCreatedAt(passage.CreatedAt),
		WithUpdatedAt(passage.UpdatedAt),
	)
}
