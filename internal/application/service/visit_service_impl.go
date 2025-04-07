package service

import (
	"blog-go/internal/domain/entity"
	"blog-go/internal/infrastructure/persistence"
	"fmt"
)

type VisitServiceImpl struct {
	db *persistence.DbHelper
}

func NewVisitServiceImpl(db *persistence.DbHelper) VisitService {
	return &VisitServiceImpl{db: db}
}

func (v *VisitServiceImpl) SaveRecord(visit entity.Visit) error {
	return v.db.VisitRepository.SaveRecord(visit)
}

func (v *VisitServiceImpl) Count(option string, parma string) (int64, error) {
	switch option {
	case "path":
		return v.db.VisitRepository.CountPath(parma)
	case "ip":
		return v.db.VisitRepository.CountIp(parma)
	case "browser":
		return v.db.VisitRepository.CountBrowser(parma)
	case "os":
		return v.db.VisitRepository.CountOs(parma)
	case "platform":
		return v.db.VisitRepository.CountPlatform(parma)
	case "all":
		return v.db.VisitRepository.CountAll()
	default:
		return 0, fmt.Errorf("Unknown option")
	}
}
