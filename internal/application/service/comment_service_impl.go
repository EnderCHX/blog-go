package service

import (
	"blog-go/internal/domain/dto"
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"
)

type CommentServiceImpl struct {
	db *persistence.DbHelper
}

func NewCommentServiceImpl(db *persistence.DbHelper) CommentService {
	return &CommentServiceImpl{db: db}
}
func (c *CommentServiceImpl) GetCommentsByPassageId(passageId string) ([]dto.CommentDTO, error) {
	comments, err := c.db.CommentRepository.GetCommentsByPassageId(passageId)
	if err != nil {
		return nil, err
	}

	return dto.CommentsToDTO(comments), nil
}

func (c *CommentServiceImpl) GetCommentsByUser(username string) ([]dto.CommentDTO, error) {
	comments, err := c.db.CommentRepository.GetCommentsByUser(username)
	if err != nil {
		return nil, err
	}

	return dto.CommentsToDTO(comments), nil
}

func (c *CommentServiceImpl) GetCommentById(commentId string) (dto.CommentDTO, error) {
	comment, err := c.db.CommentRepository.GetCommentById(commentId)
	if err != nil {
		return dto.CommentDTO{}, err
	}
	return dto.CommentToDTO(comment), nil
}

func (c *CommentServiceImpl) AddComment(commentdto dto.CommentAddDTO) (string, error) {
	comment := commentdto.ToComment()
	err := c.db.CommentRepository.AddComment(*comment)
	if err != nil {
		return "", err
	}
	return comment.CommentId, nil
}

func (c *CommentServiceImpl) DeleteComment(commentId string) error {
	return c.db.CommentRepository.DeleteComment(entity.Comment{CommentId: commentId})
}
