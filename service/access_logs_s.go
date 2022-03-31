// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"ohurlshortener/core"
	"ohurlshortener/storage"
	"ohurlshortener/utils"
	"time"
)

const access_logs_prefix = "OH_ACCESS_LOGS#"

func NewAccessLog(url string, ip string, useragent string, referer string) error {

	l := core.AccessLog{
		ShortUrl:   url,
		AccessTime: time.Now(),
		Ip:         sql.NullString{String: ip, Valid: true},
		UserAgent:  sql.NullString{String: useragent, Valid: true},
	}

	logJson, _ := json.Marshal(l)
	key := fmt.Sprintf("%s%s", access_logs_prefix, utils.UserAgentIpHash(useragent, ip))
	err := storage.RedisSet30m(key, logJson)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	return nil
}

func StoreAccessLogs() error {
	keys, err := storage.RedisScan4Keys(access_logs_prefix + "*")
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	logs := []core.AccessLog{}
	for _, k := range keys {
		v, err := storage.RedisGetString(k)
		if err != nil {
			log.Printf("redis error for key %s", k)
			continue
		}
		log := core.AccessLog{}
		json.Unmarshal([]byte(v), &log)
		logs = append(logs, log)
	} //end of for

	err = storage.InsertAccessLogs(logs)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	err = storage.RedisDelete(keys...)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}
	return nil
}

func GetPagedAccessLogs(url string, page int, size int) ([]core.AccessLog, error) {
	if page < 1 || size < 1 {
		return nil, nil
	}
	allAccessLogs, err := storage.FindAllAccessLogs(url, page, size)
	if err != nil {
		log.Println(err)
		return allAccessLogs, utils.RaiseError("内部错误，请联系管理员")
	}
	return allAccessLogs, nil
}
