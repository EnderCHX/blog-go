package repository

import "blog-go/internal/domain/entity"

type CommentRepository interface {
	GetCommentsByPassageId(passageId string) ([]entity.Comment, error)
	GetCommentsByUser(username string) ([]entity.Comment, error)
	GetCommentById(commentId string) (entity.Comment, error)
	AddComment(comment entity.Comment) error
	DeleteComment(comment entity.Comment) error
}
