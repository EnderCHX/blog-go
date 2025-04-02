package entity

type Visit struct {
	Path      string `json:"path"`
	FullPath  string `json:"full_path"`
	Ip        string `json:"ip"`
	Referer   string `json:"referer"`
	Username  string `json:"username"`
	UserAgent string `json:"user_agent"`
	Browser   string `json:"browser"`
	Platform  string `json:"platform"`
	OS        string `json:"os"`
}
