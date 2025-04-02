package midware

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

func CountVisitor(db *persistence.DbHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			visitor := entity.Visit{}
			visitor.Path = c.Request.URL.Path
			visitor.FullPath = c.FullPath()
			visitor.Ip = c.ClientIP()
			visitor.UserAgent = c.Request.UserAgent()
			ua := useragent.New(c.Request.UserAgent())
			visitor.Browser, _ = ua.Browser()
			visitor.OS = ua.OS()
			if ua.Mobile() {
				visitor.Platform = "Mobile"
			} else {
				visitor.Platform = "Desktop"
			}

			db.VisitRepository.SaveRecord(visitor)

		}()

		c.Next()
	}
}
