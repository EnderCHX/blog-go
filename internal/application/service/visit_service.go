package service

import "blog-go/internal/domain/entity"

type VisitService interface {
	SaveRecord(visit entity.Visit) error
	Count(option string, parma string) (int64, error)
}
