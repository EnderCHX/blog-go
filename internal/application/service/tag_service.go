package service

type TagService interface {
	GetTags() ([]string, error)
	GetPassageTags(passageId string) ([]string, error)
	GetTagName(id int) (string, error)
	GetTagId(tagname string) (int, error)
	AddTags(tags []string) error
	DeleteTags(tags []string) error
}
