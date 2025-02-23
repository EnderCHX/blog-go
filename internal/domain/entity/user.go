package entity

type User struct {
	Username string `json:"username" gorm:"primaryKey;index;not null"`
}
