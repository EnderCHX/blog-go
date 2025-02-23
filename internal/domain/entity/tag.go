package entity

type Tag struct {
	TagId   int    `json:"tag_id" gorm:"primaryKey;autoIncrement"`
	TagName string `json:"tag_name" gorm:"index;unique;not null"`
}

type TagOptions func(*Tag)

func NewTag(option ...TagOptions) *Tag {
	tag := &Tag{}
	for _, opt := range option {
		opt(tag)
	}
	return tag
}

func WithTagId(tagId int) TagOptions {
	return func(tag *Tag) {
		tag.TagId = tagId
	}
}

func WithTagName(tagName string) TagOptions {
	return func(tag *Tag) {
		tag.TagName = tagName
	}
}
