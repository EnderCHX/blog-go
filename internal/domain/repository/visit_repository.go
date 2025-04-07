package repository

import "blog-go/internal/domain/entity"

type VisitRepository interface {
	SaveRecord(visit entity.Visit) error
	CountAll() (int64, error)
	CountPath(path string) (int64, error)
	CountIp(ip string) (int64, error)
	CountBrowser(browser string) (int64, error)
	CountOs(os string) (int64, error)
	CountPlatform(platform string) (int64, error)
}
