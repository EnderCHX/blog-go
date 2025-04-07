package handle

import (
	"blog-go/internal/application/service"
	"blog-go/internal/domain/entity"
	"blog-go/internal/interfaces/response"
	"github.com/gin-gonic/gin"
)

type VisitHandle struct {
	VisitService service.VisitService
}

func NewVisitHandle(visitService service.VisitService) *VisitHandle {
	return &VisitHandle{
		VisitService: visitService,
	}
}

// "/visit"
func (h *VisitHandle) GetVisitCount(c *gin.Context) {
	option := c.Query("option")
	parma := c.Query("param")

	if option != "path" || parma != "/visitor_count_github" {
		_, ok := c.Get("username")
		role, _ := c.Get("role")
		if !ok {
			c.JSON(response.Unauthorized(nil))
			return
		}
		if role != entity.ADMIN {
			c.JSON(response.Default(200, false, response.USER_PERMISSION_DENIED, "权限不足", nil))
			return
		}
	}

	count, err := h.VisitService.Count(option, parma)
	if err != nil {
		c.JSON(response.InternalServerError(nil))
		return
	}
	c.JSON(response.Success(gin.H{
		"count": count,
	}))
}
