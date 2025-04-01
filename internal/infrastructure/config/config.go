package config

import (
	"encoding/json"
	"log"
	"os"
)

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type MySQL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type ApiConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Mode       string `json:"mode"`
	UserApiUrl string `json:"user_api_url"`
}

type LogCongfig struct {
	LogPath  string `json:"log_path"`
	LogLevel string `json:"log_level"`
}

type SecretKeys struct {
	RefreshSecret string `json:"refresh_secret"`
	AccessSecret  string `json:"access_secret"`
	PasswdSecret  string `json:"passwd_secret"`
}

type MailConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	MySQLConfig MySQL      `json:"mysql_config"`
	RedisConfig Redis      `json:"redis_config"`
	ApiConfig   ApiConfig  `json:"api_config"`
	LogCongfig  LogCongfig `json:"log_config"`
	SecretKeys  SecretKeys `json:"secret_keys"`
	MailConfig  MailConfig `json:"mail_config"`
}

var ConfigContext Config

var DefaultConfig = Config{
	MySQLConfig: MySQL{
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "root",
		DBName:   "blog",
	},
	RedisConfig: Redis{
		Host:     "127.0.0.1",
		Port:     "6379",
		Username: "",
		Password: "",
		DB:       0,
	},
	ApiConfig: ApiConfig{
		Host:       "0.0.0.0",
		Port:       "1314",
		Mode:       "release",
		UserApiUrl: "https://api.passport.hrbeu.top",
	},
	LogCongfig: LogCongfig{
		LogLevel: "debug",
		LogPath:  "./log.log",
	},
	SecretKeys: SecretKeys{
		RefreshSecret: "refresh_secret",
		AccessSecret:  "access_secret",
		PasswdSecret:  "passwd_secret",
	},
	MailConfig: MailConfig{
		Host:     "smtp.qq.com",
		Port:     "465",
		Username: "xxxxxxxx@qq.com",
		Password: "xxxxxxxxxxxxxxxx",
	},
}

func Setup(ConfigFileName string) {
	log.Println("读取配置文件")
	ConfigFile, err := os.ReadFile(ConfigFileName)
	if err != nil {
		log.Println("配置文件不存在，使用默认配置")
		ConfigContext = DefaultConfig
		data, _ := json.Marshal(DefaultConfig)
		os.WriteFile(ConfigFileName, data, 0644)
		return
	}
	err = json.Unmarshal(ConfigFile, &ConfigContext)
	if err != nil {
		log.Panicln(err)
		return
	}
	log.Println("读取配置文件成功")
	// fmt.Println(ConfigContext)
}

func GetConfig() Config { return ConfigContext }
