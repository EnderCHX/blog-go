package midware

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/config"
	"blog-go/internal/infrastructure/persistence"
	"github.com/EnderCHX/chx-tools-go/auth"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

func CountVisitor(db *persistence.DbHelper, cf config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			visitor := entity.Visit{}
			visitor.Path = c.Request.URL.Path
			visitor.FullPath = c.FullPath()
			visitor.Referer = c.Request.Referer()
			visitor.Ip = c.ClientIP()
			visitor.UserAgent = c.Request.UserAgent()
			ua := useragent.New(c.Request.UserAgent())
			visitor.Browser, _ = ua.Browser()
			visitor.OS = ua.OS()
			visitor.Time = time.Now()
			if ua.Mobile() {
				visitor.Platform = "Mobile"
			} else {
				visitor.Platform = "Desktop"
			}
			token := c.GetHeader("Authorization")
			if token != "" {
				token = strings.Replace(token, "Bearer ", "", 1)
				payload, err := auth.VerifyToken(token, cf.SecretKeys.AccessSecret)
				if err == nil {
					visitor.Username = payload.Username
				} else {
					visitor.Username = "GUEST"
				}
			}

			db.VisitRepository.SaveRecord(visitor)

		}()

		c.Next()
	}
}
