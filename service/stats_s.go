// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package service

import (
	"log"
	"ohurlshortener/core"
	"ohurlshortener/storage"
	"ohurlshortener/utils"
)

func GetShortUrlStats(url string) (core.ShortUrlStats, error) {
	found, err := storage.GetUrlStats(url)
	if err != nil {
		log.Println(err)
		return found, utils.RaiseError("内部错误，请联系管理员")
	}
	return found, nil
}
