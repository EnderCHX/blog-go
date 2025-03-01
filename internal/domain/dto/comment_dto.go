package dto

import (
	"blog-go/internal/domain/entity"
	"time"
)

type CommentDTO struct {
	CommentId      string    `json:"comment_id"`
	ReplyCommentId string    `json:"reply_comment_id"`
	PassageId      string    `json:"passage_id"`
	Username       string    `json:"username"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

func CommentToDTO(comment entity.Comment) CommentDTO {
	return CommentDTO{
		CommentId:      comment.CommentId,
		ReplyCommentId: comment.ReplyCommentId,
		PassageId:      comment.PassageId,
		Username:       comment.Username,
		Content:        comment.Content,
		CreatedAt:      comment.CreatedAt,
	}
}

func CommentsToDTO(comments []entity.Comment) []CommentDTO {
	var commentDTOs []CommentDTO
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, CommentToDTO(comment))
	}
	return commentDTOs
}

type CommentAddDTO struct {
	PassageId      string `json:"passage_id"`
	Content        string `json:"content"`
	ReplyCommentId string `json:"reply_comment_id"`
	Username       string `json:"username"`
}

func (c *CommentAddDTO) ToComment() *entity.Comment {
	return entity.NewComment(
		entity.WithCommentContent(c.Content),
		entity.WithCommentPassageId(c.PassageId),
		entity.WithReplyCommentId(c.ReplyCommentId),
		entity.WithCommentCreatedAt(time.Now()),
		entity.WithUsername(c.Username),
	).GenCommentId()
}
