package handle

import (
	"blog-go/internal/application/service"
	"blog-go/internal/domain/dto"
	"blog-go/internal/domain/entity"
	"blog-go/internal/interfaces/response"
	"errors"
	"github.com/EnderCHX/chx-tools-go/auth"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PassageHandle struct {
	passageService service.PassageService
	logger         *zap.Logger
}

func NewPassageHandle(passageService service.PassageService, logger *zap.Logger) *PassageHandle {
	return &PassageHandle{
		passageService: passageService,
		logger:         logger,
	}
}

// "/post/:id"
func (h *PassageHandle) GetPassage(c *gin.Context) {
	passageId := c.Param("id")
	passage, err := h.passageService.GetPassageById(passageId)

	if err != nil {
		h.logger.Error("[PassageService] 获取文章 "+passageId+" 失败", zap.String("error", err.Error()))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(response.NotFound(nil))
			return
		}

		c.JSON(response.InternalServerError(nil))
		return
	}

	c.JSON(response.Success(dto.PassageToDTO(passage)))
}

// "/posts"
func (h *PassageHandle) GetPassages(c *gin.Context) {
	passages, err := h.passageService.GetPassages()

	if err != nil {
		h.logger.Error("[PassageService] 获取文章列表失败", zap.String("error", err.Error()))
		c.JSON(response.InternalServerError(nil))
		return
	}

	c.JSON(response.Success(dto.PassagesToDTO(passages)))
}

// "/posts/tag/:tagname"
func (h *PassageHandle) GetPassagesByTag(c *gin.Context) {
	tagname := c.Param("tagname")
	passages, err := h.passageService.GetPassagesByTag(tagname)

	if err != nil {
		h.logger.Error("[PassageService] 获取文章列表失败", zap.String("error", err.Error()))
		c.JSON(response.InternalServerError(nil))
	}

	c.JSON(response.Success(passages))
}

// "posts/date/:start/:end"
func (h *PassageHandle) GetPassagesByDate(c *gin.Context) {
	start, err := time.Parse("2006-01-02", c.Param("start"))
	if err != nil {
		c.JSON(response.BadRequest(nil))
		return
	}

	end, err := time.Parse("2006-01-02", c.Param("end"))
	if err != nil {
		c.JSON(response.BadRequest(nil))
		return
	}

	passages, err := h.passageService.GetPassagesByDate(start, end)
	if err != nil {
		h.logger.Error("[PassageService] 获取文章列表失败", zap.String("error", err.Error()))
		c.JSON(response.InternalServerError(nil))
	}

	c.JSON(response.Success(passages))
}

func (h *PassageHandle) AddPassage(c *gin.Context) {
	claims, _ := c.Get("claims")
	if claims == nil {
		c.JSON(response.Unauthorized(nil))
		return
	}

	var passageDTO struct {
		dto.PassageDTO
		dto.TagsDTO
	}
	c.BindJSON(&passageDTO)

	passage := *entity.NewPassage(
		entity.WithTitle(passageDTO.Title),
		entity.WithContent(passageDTO.Content),
		entity.WithAuthorUsername(claims.(*auth.JWTPayload).Username),
		entity.WithCreatedAt(time.Now()),
		entity.WithUpdatedAt(time.Now()),
	).GenPassageId()

	err := h.passageService.AddPassage(passage)

	if err != nil {
		c.JSON(response.InternalServerError(nil))
		return
	}

	c.JSON(response.Success(gin.H{
		"passage_id": passage.GenPassageId().PassageId,
		"title":      passage.Title,
		"author":     passage.AuthorUsername,
		"tags":       passageDTO.Tags,
	}))

	err = h.passageService.AddTags(passageDTO.Tags)
	if err != nil {
		h.logger.Warn("[PassageService] 创建标签失败", zap.String("error", err.Error()))
	}
	err = h.passageService.AddPassageTags(passage.PassageId, passageDTO.Tags)
	if err != nil {
		h.logger.Warn("[PassageService] 添加标签失败", zap.String("error", err.Error()))
	}
}
