package service

import (
	"blog-go/internal/domain/dto"
)

type CommentService interface {
	GetCommentsByPassageId(passageId string) ([]dto.CommentDTO, error)
	GetCommentsByUser(username string) ([]dto.CommentDTO, error)
	GetCommentById(commentId string) (dto.CommentDTO, error)
	AddComment(commentdto dto.CommentAddDTO) (string, error)
	DeleteComment(commentId string) error
}
