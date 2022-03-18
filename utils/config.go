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
	Port      int
	AdminPort int
	UrlPrefix string
	Debug     bool
}

type RedisConfigInfo struct {
	Host     string
	User     string
	Password string
	Database int
	PoolSize int
}

type DatabaseConfigInfo struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	MaxOpenConns int
	MaxIdleConn  int
}

func InitConfig(file string) (*ini.File, error) {

	cfg, err := ini.Load(file)
	if err != nil {
		return nil, nil
	}

	section := cfg.Section("postgres")
	DatabaseConifg.Host = section.Key("host").String()
	DatabaseConifg.Port = section.Key("port").MustInt()
	DatabaseConifg.MaxOpenConns = section.Key("max_open_conn").MustInt()
	DatabaseConifg.MaxIdleConn = section.Key("max_idle_conn").MustInt()
	DatabaseConifg.User = section.Key("user").String()
	DatabaseConifg.Password = section.Key("password").String()
	DatabaseConifg.DbName = section.Key("database").String()

	appSection := cfg.Section("app")
	AppConfig.Debug = appSection.Key("debug").MustBool()
	AppConfig.Port = appSection.Key("port").MustInt()
	AppConfig.AdminPort = appSection.Key("admin_port").MustInt()
	AppConfig.UrlPrefix = appSection.Key("url_prefix").String()

	redisSection := cfg.Section("redis")
	RedisConfig.Host = redisSection.Key("host").String()
	RedisConfig.User = redisSection.Key("user").String()
	RedisConfig.Password = redisSection.Key("password").String()
	RedisConfig.Database = redisSection.Key("database").MustInt()
	RedisConfig.PoolSize = redisSection.Key("pool_size").MustInt()

	return cfg, err
}
