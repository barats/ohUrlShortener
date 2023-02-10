// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package service

import (
	"ohurlshortener/core"
	"ohurlshortener/storage"
	"ohurlshortener/utils"
)

// GetSumOfUrlStats 获取所有短链接的统计信息
func GetSumOfUrlStats() (int, core.ShortUrlStats, error) {
	var (
		totalCount int
		result     core.ShortUrlStats
	)

	totalCount, err := storage.GetUrlCount()
	if err != nil {
		return totalCount, result, utils.RaiseError("内部错误，请联系管理员！")
	}

	result, er := storage.GetSumOfUrlStats()
	if er != nil {
		return totalCount, result, utils.RaiseError("内部错误，请联系管理员！")
	}

	return totalCount, result, nil
}

// GetShortUrlStats 获取单个短链接的统计信息
func GetShortUrlStats(url string) (core.ShortUrlStats, error) {
	found, err := storage.GetUrlStats(url)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员！")
	}
	return found, nil
}

// GetTop25Url 获取访问量最高的 25 个短链接
func GetTop25Url() ([]core.Top25Url, error) {
	found, err := storage.GetTop25()
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员！")
	}
	return found, nil
}

// GetPagedUrlIpCountStats 获取单个短链接的 IP 访问量统计信息
func GetPagedUrlIpCountStats(url string, page int, size int) ([]core.UrlIpCountStats, error) {
	if page < 1 || size < 1 {
		return nil, nil
	}
	found, err := storage.FindPagedUrlIpCountStats(url, page, size)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员！")
	}
	return found, nil
}
