package blog

import "blog-go/database"

type PassageTags struct {
	PassageId int `json:"passage_id" gorm:"primaryKey"`
	TagId     int `json:"tag_id" gorm:"not null"`
}

func CreatePassageTags() {
	db := database.GetDB()
	db.AutoMigrate(&PassageTags{})
}
