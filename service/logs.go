package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"ohurlshortener/core"
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"
	"time"
)

const access_logs_prefix = "OH_ACCESS_LOGS#"

func NewAccessLog(url string, ip string, useragent string) error {

	l := core.AccessLog{
		ShortUrl:   url,
		AccessTime: time.Now(),
		Ip:         sql.NullString{String: ip, Valid: true},
		UserAgent:  sql.NullString{String: useragent, Valid: true},
	}

	logJson, _ := json.Marshal(l)
	key := fmt.Sprintf("%s%s", access_logs_prefix, utils.UserAgentIpHash(useragent, ip))
	err := redis.Set30m(key, logJson)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	return nil
}

func StoreAccessLog() error {
	keys, err := redis.Scan4Keys(access_logs_prefix + "*")
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	logs := []core.AccessLog{}
	for _, k := range keys {
		v, err := redis.GetString(k)
		if err != nil {
			log.Printf("redis error for key %s", k)
			continue
		}
		log := core.AccessLog{}
		json.Unmarshal([]byte(v), &log)
		logs = append(logs, log)
	} //end of for

	err = db.InsertAccessLogs(logs)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	err = redis.Delete(keys...)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}
	return nil
}
