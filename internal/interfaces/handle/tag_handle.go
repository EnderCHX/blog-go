package handle

import (
	"blog-go/internal/application/service"
	"blog-go/internal/interfaces/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TagHandle struct {
	tagService service.TagService
	logger     *zap.Logger
}

func NewTagHandle(tagService service.TagService, logger *zap.Logger) *TagHandle {
	return &TagHandle{
		tagService: tagService,
		logger:     logger,
	}
}

func (h *TagHandle) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetTags()
	if err != nil {
		h.logger.Error("GetTags error", zap.Error(err))
		c.JSON(response.InternalServerError(nil))
		return
	}
	c.JSON(response.Success(tags))
}

// "post/:id/tags"
func (h *TagHandle) GetPassageTags(c *gin.Context) {
	passageId := c.Param("id")
	tags, err := h.tagService.GetPassageTags(passageId)
	if err != nil {
		h.logger.Error("GetPassageTags error", zap.Error(err))
		c.JSON(response.InternalServerError(nil))
		return
	}
	c.JSON(response.Success(tags))
}
