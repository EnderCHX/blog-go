package handle

import (
	"blog-go/internal/application/service"
	"blog-go/internal/domain/dto"
	"blog-go/internal/interfaces/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentHandle struct {
	CommentService service.CommentService
	UserService    service.UserService
	logger         *zap.Logger
}

func NewCommentHandle(commentService service.CommentService, userService service.UserService, logger *zap.Logger) *CommentHandle {
	return &CommentHandle{
		CommentService: commentService,
		UserService:    userService,
		logger:         logger,
	}
}

// "/comment/passage/:passageId"
func (c *CommentHandle) GetCommentsByPassageId(ctx *gin.Context) {
	passageId := ctx.Param("passageId")
	comments, err := c.CommentService.GetCommentsByPassageId(passageId)
	if err != nil {
		ctx.JSON(response.InternalServerError(nil))
		return
	}

	ctx.JSON(response.Success(comments))
}

// "/comments/user/:username"
func (c *CommentHandle) GetCommentsByUser(ctx *gin.Context) {
	username := ctx.Param("username")
	comments, err := c.CommentService.GetCommentsByUser(username)
	if err != nil {
		ctx.JSON(response.InternalServerError(nil))
		return
	}

	ctx.JSON(response.Success(comments))
}

// "/comment/:commentId"
func (c *CommentHandle) GetCommentById(ctx *gin.Context) {
	commentId := ctx.Param("commentId")
	comment, err := c.CommentService.GetCommentById(commentId)
	if err != nil {
		ctx.JSON(response.InternalServerError(nil))
		return
	}

	ctx.JSON(response.Success(comment))
}

// "/newcomment"
func (c *CommentHandle) AddComment(ctx *gin.Context) {
	username, ok := ctx.Get("username")
	if !ok {
		ctx.JSON(response.Unauthorized(nil))
		return
	}

	go func() {
		token := ctx.GetHeader("Authorization")
		token = token[7:]
		err := c.UserService.SavaEmail(username.(string), token)
		if err != nil {
			c.logger.Error("SavaEmail error", zap.Error(err))
		}
	}()

	var comment dto.CommentAddDTO
	ctx.BindJSON(&comment)
	comment.Username = username.(string)

	commentId, err := c.CommentService.AddComment(comment)
	if err != nil {
		ctx.JSON(response.InternalServerError(nil))
		return
	}

	go func() {
		replyComment, err := c.CommentService.GetCommentById(comment.ReplyCommentId)
		if err != nil {
			c.logger.Error("GetCommentById error", zap.Error(err))
			return
		}
		if replyComment.Username != "" {
			err = c.UserService.SendEmail(replyComment.Username, "回复通知",
				"<p>您收到一条回复通知，请及时查看。</p>"+
					"<p>回复人："+comment.Username+"</p>"+
					"<p>"+comment.Content+"</p>",
			)
			if err != nil {
				c.logger.Error("SendEmail error", zap.Error(err))
				return
			}
		}
	}()

	c.logger.Info("AddComment", zap.String("comment", comment.Content))
	ctx.JSON(response.Success(gin.H{
		"comment_id": commentId,
		"":           comment,
	}))
}

func (c *CommentHandle) DeleteComment(ctx *gin.Context) {
	ctx.JSON(response.Unauthorized(nil))
}
