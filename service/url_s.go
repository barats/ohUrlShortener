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
	"ohurlshortener/storage"
	"ohurlshortener/utils"
	"time"
)

// 从数据库中获取所有「有效」状态的短链接
// 并将其可以 key-> value 形式存入 Redis 中
func ReloadUrls() (bool, error) {
	//把所有访问日志记录到数据库中
	err := StoreAccessLogs()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	//找出所有已经配置好的短链接
	urls, err := storage.FindAllShortUrls()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	//清理 redis db
	err = storage.RedisFlushDB()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	//将所有「有效」状态的短域名再次放入 Redis
	for _, url := range urls {
		if url.Valid {
			err := storage.RedisSet4Ever(url.ShortUrl, url.DestUrl)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	} //end of for
	return true, nil
}

// 从 Redis 中查询目标短链接是否存在
func Search4ShortUrl(shortUrl string) (string, error) {
	destUrl, err := storage.RedisGetString(shortUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}
	return destUrl, nil
}

// 获取分页的短链接信息
func GetPagesShortUrls(url string, page int, size int) ([]core.ShortUrl, error) {
	if page < 1 || size < 1 {
		return nil, nil
	}
	allUrls, err := storage.FindPagedShortUrls(url, page, size)
	if err != nil {
		log.Println(err)
		return allUrls, utils.RaiseError("内部错误，请联系管理员")
	}
	return allUrls, nil
}

// 生成短链接
func GenerateShortUrl(destUrl string, memo string) (string, error) {
	shortUrl, err := core.GenerateShortLink(destUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	foundUrl, err := storage.FindShortUrl(shortUrl)
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

	if err := storage.InsertShortUrl(url); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	if err := storage.RedisSet4Ever(shortUrl, destUrl); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	return shortUrl, nil
}

//禁用/启用短链接
func ChangeState(shortUrl string, enable bool) (bool, error) {
	found, err := storage.FindShortUrl(shortUrl)
	if err != nil {
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	if found.IsEmpty() {
		return false, utils.RaiseError("该短链接不存在")
	}

	found.Valid = enable

	e := storage.UpdateShortUrl(found)
	if e != nil {
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	if enable {
		storage.RedisSet4Ever(found.ShortUrl, found.DestUrl)
	} else {
		storage.RedisDelete(found.ShortUrl)
	}

	return true, nil
}
