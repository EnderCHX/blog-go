package handle

import (
	"blog-go/internal/application/service"
	"blog-go/internal/domain/dto"
	"blog-go/internal/interfaces/response"
	"github.com/EnderCHX/chx-tools-go/auth"
	"github.com/gin-gonic/gin"
)

type CommentHandle struct {
	CommentService service.CommentService
}

func NewCommentHandle(commentService service.CommentService) *CommentHandle {
	return &CommentHandle{
		CommentService: commentService,
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
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(response.Unauthorized(nil))
		return
	}

	var comment dto.CommentAddDTO
	ctx.BindJSON(&comment)
	comment.Username = claims.(*auth.JWTPayload).Username

	commentId, err := c.CommentService.AddComment(comment)
	if err != nil {
		ctx.JSON(response.InternalServerError(nil))
		return
	}

	ctx.JSON(response.Success(gin.H{
		"comment_id": commentId,
		"":           comment,
	}))
}

func (c *CommentHandle) DeleteComment(ctx *gin.Context) {
	ctx.JSON(response.Unauthorized(nil))
}
