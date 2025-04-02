package entity

import "sync"

type PathCount struct {
	FullPath string `json:"full_path" gorm:"primaryKey"`
	Count    int    `json:"count"`
	Lock     sync.Mutex
}

type IpCount struct {
	Ip       string `json:"ip" gorm:"primaryKey"`
	Count    int    `json:"count"`
	FullPath string `json:"full_path"`
	Lock     sync.Mutex
}
