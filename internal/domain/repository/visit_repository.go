package repository

import "blog-go/internal/domain/entity"

type VisitRepository interface {
	SaveRecord(visit entity.Visit) error
}
