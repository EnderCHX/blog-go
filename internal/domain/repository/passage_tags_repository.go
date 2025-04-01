package repository

type PassageTagsRepository interface {
	GetPassageTags(passageId string) ([]string, error)
	AddPassageTags(passageId string, tags []string) error
	GetTagPassages(tagName string) ([]string, error)
}
