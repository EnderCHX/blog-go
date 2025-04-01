package persistence

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/cache"
	"blog-go/internal/infrastructure/utils"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CommentRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.Redis
}

func (c *CommentRepositoryImpl) AutoMigrate() {
	c.db.AutoMigrate(&entity.Comment{})
}

func (c *CommentRepositoryImpl) SetDb(db *gorm.DB) {
	c.db = db
}

func (c *CommentRepositoryImpl) SetRedis(rdb *cache.Redis) {
	c.rdb = rdb
}

func (c *CommentRepositoryImpl) GetCommentsByPassageId(passageId string) ([]entity.Comment, error) {
	var comments []entity.Comment
	commentsStr, err := c.rdb.Get("blog:comments:passage:" + passageId)
	if err != nil {
		err2 := c.db.Model(&entity.Comment{}).
			Where("passage_id = ?", passageId).
			Where("deleted = ?", 0).
			Order("created_at desc").
			Find(&comments).Error
		if err2 != nil {
			return nil, err2
		}

		if err == redis.Nil {
			go func() {
				commentsToJson, _ := json.Marshal(comments)
				c.rdb.Set("blog:comments:passage:"+passageId, string(commentsToJson), utils.RandExpTime(60, 70))
			}()
		}

		return comments, nil
	}

	err = json.Unmarshal([]byte(commentsStr), &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentRepositoryImpl) GetCommentsByUser(username string) ([]entity.Comment, error) {
	var comments []entity.Comment
	commentsStr, err := c.rdb.Get("blog:comments:user:" + username)
	if err != nil {
		err2 := c.db.Model(&entity.Comment{}).
			Where("username = ?", username).
			Where("deleted = ?", 0).
			Order("created_at desc").
			Find(&comments).Error
		if err2 != nil {
			return nil, err2
		}

		if err == redis.Nil {
			go func() {
				commentsToJson, _ := json.Marshal(comments)
				c.rdb.Set("blog:comments:user:"+username, string(commentsToJson), utils.RandExpTime(60, 70))
			}()
		}

		return comments, nil
	}

	err = json.Unmarshal([]byte(commentsStr), &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentRepositoryImpl) GetCommentById(commentId string) (entity.Comment, error) {
	var comment entity.Comment
	commentStr, err := c.rdb.Get("blog:comment:" + commentId)
	if err != nil {
		err2 := c.db.Model(&entity.Comment{}).
			Where("comment_id = ?", commentId).
			Where("deleted = ?", 0).
			First(&comment).Error

		if err2 != nil {
			return comment, err2
		}

		if err == redis.Nil {
			go func() {
				c.rdb.Set("blog:comment:"+commentId, comment.ToJson(), utils.RandExpTime(60, 70))
			}()
		}

		return comment, nil
	}

	err = json.Unmarshal([]byte(commentStr), &comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (c *CommentRepositoryImpl) AddComment(comment entity.Comment) error {
	err := c.db.Create(&comment).Error
	if err != nil {
		return err
	}

	err = c.rdb.Del("blog:comments:passage:" + comment.PassageId)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentRepositoryImpl) DeleteComment(comment entity.Comment) error {
	err := c.db.Model(&entity.Comment{}).
		Where("comment_id = ?", comment.CommentId).
		Update("deleted", 1).Error

	if err != nil {
		return err
	}

	var passageid string
	err = c.db.Model(&entity.Comment{}).
		Where("comment_id = ?", comment.CommentId).
		Pluck("passage_id", &passageid).Error
	if err != nil {
		return err
	}

	err = c.rdb.Del("blog:comments:passage:" + passageid)
	if err != nil {
		return err
	}
	return nil
}
