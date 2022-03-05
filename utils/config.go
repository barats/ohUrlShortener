package utils

import (
	"gopkg.in/ini.v1"
)

var (
	DatabaseConifg DatabaseConfigInfo
	AppConfig      AppConfigInfo
	RedisConfig    RedisConfigInfo
)

type AppConfigInfo struct {
	Port int
}

type RedisConfigInfo struct {
	Host     string
	User     string
	Password string
	Database int
}

type DatabaseConfigInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func InitConfig(file string) (*ini.File, error) {

	cfg, err := ini.Load(file)
	if err != nil {
		return nil, nil
	}

	section := cfg.Section("postgres")
	DatabaseConifg.Host = section.Key("host").String()
	DatabaseConifg.Port = section.Key("port").MustInt()
	DatabaseConifg.User = section.Key("user").String()
	DatabaseConifg.Password = section.Key("password").String()
	DatabaseConifg.DbName = section.Key("database").String()

	appSection := cfg.Section("app")
	AppConfig.Port = appSection.Key("port").MustInt()

	redisSection := cfg.Section("redis")
	RedisConfig.Host = redisSection.Key("host").String()
	RedisConfig.User = redisSection.Key("user").String()
	RedisConfig.Password = redisSection.Key("password").String()
	RedisConfig.Database = redisSection.Key("database").MustInt()

	return cfg, err
}
