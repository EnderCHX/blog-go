package entity

import (
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	Username string `json:"username" gorm:"primaryKey;index;not null"`
	Email    string `json:"email" gorm:"not null"`
}

func NewUser(username, email string) *User {
	return &User{
		Username: username,
		Email:    email,
	}
}

// 角色
const (
	ADMIN  = "ADMIN"
	EDITOR = "EDITOR"
	AUTHOR = "AUTHOR"
	USER   = "USER"
	GUEST  = "GUEST"
)

func (u *User) GetEmail(apiUrl, accessToken string) (string, error) {
	req, err := http.NewRequest("GET", apiUrl+"/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		return "", nil
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()

	var data struct {
		Data struct {
			Email string `json:"email"`
		} `json:"data"`
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", nil
	}

	return data.Data.Email, nil
}
