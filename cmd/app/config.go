package main

import "os"

func GetAppEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	return env
}

type AppConfig struct {
	Env         string       `json:"env"`
	Host        string       `json:"host"`
	Port        string       `json:"port"`
	Redis       RedisConfig  `json:"redisConfig"`
	TgiServices []TgiService `json:"tgiServices"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type TgiService struct {
	Scheme    string `json:"scheme"`
	Host      string `json:"host"`
	Port      int64  `json:"port"`
	Path      string `json:"path"`
	Token     string `json:"token"`
	ServiceId string `json:"serviceId"`
}
