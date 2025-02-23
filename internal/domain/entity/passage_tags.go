package entity

type PassageTags struct {
	PassageId string `json:"passage_id" gorm:"not null"`
	TagId     int    `json:"tag_id" gorm:"not null"`
}

type PassageTagsOptions func(*PassageTags)

func NewPassageTags(option ...PassageTagsOptions) *PassageTags {
	passageTags := &PassageTags{}
	for _, opt := range option {
		opt(passageTags)
	}
	return passageTags
}
