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
	"log"
	"ohurlshortener/core"
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"
	"time"
)

func ReloadUrls() (bool, error) {
	//把所有访问日志记录到数据库中
	err := StoreAccessLogs()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	//找出所有已经配置好的短链接
	urls, err := db.FindAllShortUrls()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	//清理 redis db
	err = redis.FlushDB()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	//将所有「有效」状态的短域名再次放入 Redis
	for _, url := range urls {
		if url.Valid {
			err := redis.Set4Ever(url.ShortUrl, url.DestUrl)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	} //end of for
	return true, nil
}

func Search4ShortUrl(shortUrl string) (string, error) {
	destUrl, err := redis.GetString(shortUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}
	return destUrl, nil
}

func GetPagesShortUrls(url string, page int, size int) ([]core.ShortUrl, error) {
	if page < 1 || size < 1 {
		return nil, nil
	}
	allUrls, err := db.FindPagedShortUrls(url, page, size)
	if err != nil {
		log.Println(err)
		return allUrls, utils.RaiseError("内部错误，请联系管理员")
	}
	return allUrls, nil
}

func GenerateShortUrl(destUrl string, memo string) (string, error) {
	shortUrl, err := core.GenerateShortLink(destUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	foundUrl, err := db.FindShortUrl(shortUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	if !foundUrl.IsEmpty() {
		return shortUrl, nil
	}

	var nsMemo sql.NullString
	if !utils.EemptyString(memo) {
		nsMemo = sql.NullString{Valid: true, String: memo}
	}

	url := core.ShortUrl{
		DestUrl:   destUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now(),
		Valid:     true,
		Memo:      nsMemo,
	}

	if err := db.InsertShortUrl(url); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	if err := redis.Set4Ever(shortUrl, destUrl); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	return shortUrl, nil
}
