package repository

type PassageTagsRepository interface {
	GetPassageTags(passageId string) ([]string, error)
	GetTagPassages(tagName string) ([]string, error)
}
