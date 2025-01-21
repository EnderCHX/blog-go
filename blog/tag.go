package blog

type Tag struct {
	TagId   int    `json:"tag_id" gorm:"primaryKey"`
	TagName string `json:"tag_name"`
}
