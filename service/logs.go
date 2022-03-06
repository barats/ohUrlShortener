package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"ohurlshortener/core"
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"
	"time"
)

const key_access_logs_list = "OH_CURRENT_ACCESS_LOGS"

func NewAccessLog(url string, ip string, useragent string) error {

	l := core.AccessLog{
		ShortUrl:   url,
		AccessTime: time.Now(),
		Ip:         sql.NullString{String: ip, Valid: true},
		UserAgent:  sql.NullString{String: useragent, Valid: true},
	}

	logJson, _ := json.Marshal(l)
	err := redis.PushToList(key_access_logs_list, logJson)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	return nil
}

func StoreAccessLog() error {
	reuslt, err := redis.GetAllFromList(key_access_logs_list)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	logs := []core.AccessLog{}
	for _, r := range reuslt {
		log := core.AccessLog{}
		json.Unmarshal([]byte(r), &log)
		logs = append(logs, log)
	}

	err = db.InsertAccessLogs(logs)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	err = redis.Expire(key_access_logs_list)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	return nil
}
