package entity

import (
	"encoding/json"
	"github.com/EnderCHX/chx-tools-go/encrypt"
	"time"
)

type Comment struct {
	CommentId      string    `json:"comment_id" gorm:"primaryKey;type:char(32);not null"`
	ReplyCommentId string    `json:"reply_comment_id"`
	PassageId      string    `json:"passage_id" gorm:"index;not null"`
	Username       string    `json:"username" gorm:"index;not null"`
	Content        string    `json:"content" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	DeletedAt      time.Time `json:"deleted_at"`
	Deleted        bool      `json:"deleted"`
}

func (c *Comment) GenCommentId() *Comment {
	c.CommentId = encrypt.Md5(c.Username + c.PassageId + c.CreatedAt.String())
	return c
}

func (c *Comment) ToJson() string {
	jsonStr, _ := json.Marshal(c)
	return string(jsonStr)
}

type CommentOptions func(*Comment)

func WithCommentId(commentId string) CommentOptions {
	return func(c *Comment) {
		c.CommentId = commentId
	}
}

func WithReplyCommentId(replyCommentId string) CommentOptions {
	return func(c *Comment) {
		c.ReplyCommentId = replyCommentId
	}
}

func WithCommentPassageId(passageId string) CommentOptions {
	return func(c *Comment) {
		c.PassageId = passageId
	}
}

func WithUsername(username string) CommentOptions {
	return func(c *Comment) {
		c.Username = username
	}
}

func WithCommentContent(content string) CommentOptions {
	return func(c *Comment) {
		c.Content = content
	}
}

func WithCommentCreatedAt(createdAt time.Time) CommentOptions {
	return func(c *Comment) {
		c.CreatedAt = createdAt
	}
}

func WithCommentDeletedAt(deletedAt time.Time) CommentOptions {
	return func(c *Comment) {
		c.DeletedAt = deletedAt
	}
}

func WithCommentDeleted(deleted bool) CommentOptions {
	return func(c *Comment) {
		c.Deleted = deleted
	}
}

func NewComment(option ...CommentOptions) *Comment {
	comment := &Comment{}
	for _, o := range option {
		o(comment)
	}
	return comment
}
