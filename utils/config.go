// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package utils

import (
	"strings"

	"gopkg.in/ini.v1"
)

const Version = "1.9"

var (
	DatabaseConfig     DatabaseConfigInfo
	AppConfig          AppConfigInfo
	RedisConfig        RedisConfigInfo
	RedisClusterConfig RedisClusterConfigInfo
	CaptchaConfig      CaptchaConfigInfo
)

type CaptchaConfigInfo struct {
	Store string
}

// AppConfigInfo 应用配置
type AppConfigInfo struct {
	Port      int
	AdminPort int
	UrlPrefix string
	Debug     bool
}

// RedisClusterConfigInfo redis配置
type RedisClusterConfigInfo struct {
	Hosts    []string
	User     string
	Password string
	PoolSize int
}

// RedisConfigInfo redis配置
type RedisConfigInfo struct {
	Host     string
	User     string
	Password string
	Database int
	PoolSize int
}

// DatabaseConfigInfo 数据库配置
type DatabaseConfigInfo struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	MaxOpenConns int
	MaxIdleConn  int
}

// InitConfig 初始化配置
func InitConfig(file string) (*ini.File, error) {
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, nil
	}

	section := cfg.Section("postgres")
	DatabaseConfig.Host = section.Key("host").String()
	DatabaseConfig.Port = section.Key("port").MustInt()
	DatabaseConfig.MaxOpenConns = section.Key("max_open_conn").MustInt()
	DatabaseConfig.MaxIdleConn = section.Key("max_idle_conn").MustInt()
	DatabaseConfig.User = section.Key("user").String()
	DatabaseConfig.Password = section.Key("password").String()
	DatabaseConfig.DbName = section.Key("database").String()

	appSection := cfg.Section("app")
	AppConfig.Debug = appSection.Key("debug").MustBool()
	AppConfig.Port = appSection.Key("port").MustInt()
	AppConfig.AdminPort = appSection.Key("admin_port").MustInt()
	AppConfig.UrlPrefix = appSection.Key("url_prefix").String()

	redisSection := cfg.Section("redis")
	RedisConfig.Host = redisSection.Key("host").String()
	RedisConfig.User = redisSection.Key("username").String()
	RedisConfig.Password = redisSection.Key("password").String()
	RedisConfig.Database = redisSection.Key("database").MustInt()
	RedisConfig.PoolSize = redisSection.Key("pool_size").MustInt()

	redisClusterSection := cfg.Section("redis-cluster")
	hosts := redisClusterSection.Key("hosts").String()
	if !EmptyString(hosts) {
		hostsArr := strings.Split(hosts, ",")
		RedisClusterConfig.Hosts = hostsArr
	}
	RedisClusterConfig.User = redisClusterSection.Key("username").String()
	RedisClusterConfig.Password = redisClusterSection.Key("password").String()
	RedisClusterConfig.PoolSize = redisClusterSection.Key("pool_size").MustInt()

	captchaSection := cfg.Section("captcha")
	CaptchaConfig.Store = captchaSection.Key("store").String()

	return cfg, err
}
