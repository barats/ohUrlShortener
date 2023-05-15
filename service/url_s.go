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
	"time"

	"ohurlshortener/core"
	"ohurlshortener/storage"
	"ohurlshortener/utils"
)

// ReloadUrls
//
// 从数据库中获取所有「有效」状态的短链接
// 并将其可以 key-> value 形式存入 Redis 中
func ReloadUrls() (bool, error) {

	//Get total count to calculate page size
	count, err := storage.GetUrlCount()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	if count > 0 {
		// query for all urls by page
		totalPageCount := (count / 100) + 1 //100 at a time

		for i := 1; i <= totalPageCount; i++ {
			urls, err := storage.FindAllShortUrlsByPage(i, 100)
			if err != nil {
				log.Println(err)
				continue
			}
			go func() {
				for _, url := range urls {
					if url.Valid {
						mu := core.MemShortUrl{DestUrl: url.DestUrl, OpenType: url.OpenType}
						res, err := json.Marshal(mu)
						if err != nil {
							log.Println(err)
							continue
						}
						err = storage.RedisSet4Ever(url.ShortUrl, res)
						if err != nil {
							log.Println(err)
							continue
						}
					}
				} // end of for
			}()
		}
	}
	return true, nil
}

// Search4ShortUrl
//
// 从 Redis 中查询目标短链接是否存在
func Search4ShortUrl(shortUrl string) (url core.MemShortUrl, err error) {
	mu := core.MemShortUrl{}
	found, err := storage.RedisGetString(shortUrl)
	if err != nil {
		log.Println(err)
		return mu, utils.RaiseError("内部错误，请联系管理员")
	}
	if utils.EmptyString(found) {
		return mu, nil
	}
	return mu, json.Unmarshal([]byte(found), &mu)
}

// GetPagesShortUrls
//
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

// GenerateShortUrl
//
// 生成短链接
func GenerateShortUrl(destUrl string, memo string, openType int) (string, error) {
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
		// Already existed
		return shortUrl, utils.RaiseError(fmt.Sprintf("短链接 %s 已存在", shortUrl))
	}

	var nsMemo sql.NullString
	if !utils.EmptyString(memo) {
		nsMemo = sql.NullString{Valid: true, String: memo}
	}

	url := core.ShortUrl{
		DestUrl:   destUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now(),
		Valid:     true,
		Memo:      nsMemo,
		OpenType:  core.OpenType(openType),
	}

	if err := storage.InsertShortUrl(url); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	mu := core.MemShortUrl{DestUrl: url.DestUrl, OpenType: url.OpenType}
	res, err := json.Marshal(mu)
	if err != nil {
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	if err := storage.RedisSet4Ever(shortUrl, res); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	return shortUrl, nil
}

// ChangeState
//
// 禁用/启用短链接
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
		mu := core.MemShortUrl{DestUrl: found.DestUrl, OpenType: found.OpenType}
		res, err := json.Marshal(mu)
		if err != nil {
			return false, utils.RaiseError("内部错误，请联系管理员")
		}
		storage.RedisSet4Ever(found.ShortUrl, res)
	} else {
		storage.RedisDelete(found.ShortUrl)
	}

	return true, nil
}

// DeleteUrlAndAccessLogs 删除短链接以及对应的访问日志
func DeleteUrlAndAccessLogs(shortUrl string) error {
	found, err := storage.FindShortUrl(shortUrl)
	if err != nil {
		return utils.RaiseError("内部错误，请联系管理员")
	}

	if found.IsEmpty() {
		return utils.RaiseError("该短链接不存在")
	}

	err = storage.DeleteShortUrlWithAccessLogs(found)
	if err != nil {
		return utils.RaiseError("内部错误，无法删除短链接")
	}

	storage.RedisDelete(found.ShortUrl)

	return nil
}
