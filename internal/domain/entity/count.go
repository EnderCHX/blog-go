package entity

type PathCount struct {
	FullPath string `json:"full_path" gorm:"primaryKey"`
	Count    int    `json:"count"`
}

type IpCount struct {
	Ip       string `json:"ip" gorm:"primaryKey"`
	Count    int    `json:"count"`
	FullPath string `json:"full_path"`
}
