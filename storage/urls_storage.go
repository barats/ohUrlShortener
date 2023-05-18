// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package storage

import (
	"fmt"

	"ohurlshortener/core"
	"ohurlshortener/utils"
)

var MaxInsertCount = 1000

// UpdateShortUrl 更新短链接
func UpdateShortUrl(shortUrl core.ShortUrl) error {
	query := `UPDATE public.short_urls SET short_url = :short_url, dest_url = :dest_url, is_valid = :is_valid, memo = :memo WHERE id = :id`
	return DbNamedExec(query, shortUrl)
}

// DeleteShortUrl 删除短链接
func DeleteShortUrl(shortUrl core.ShortUrl) error {
	query := `DELETE from public.short_urls WHERE short_url = :short_url`
	return DbNamedExec(query, shortUrl)
}

// DeleteShortUrlWithAccessLogs 删除短链接以及其访问日志
func DeleteShortUrlWithAccessLogs(shortUrl core.ShortUrl) error {
	query1 := fmt.Sprintf(`DELETE from public.short_urls WHERE short_url = '%s'`, shortUrl.ShortUrl)
	query2 := fmt.Sprintf(`DELETE from public.access_logs WHERE short_url = '%s'`, shortUrl.ShortUrl)
	return DbExecTx(query1, query2)
} // end of Transaction Action

// FindShortUrl 根据短链接查找短链接信息
func FindShortUrl(url string) (core.ShortUrl, error) {
	found := core.ShortUrl{}
	query := `SELECT * FROM public.short_urls WHERE short_url = $1`
	err := DbGet(query, &found, url)
	return found, err
}

// FindAllShortUrls 查找所有短链接
func FindAllShortUrls() ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	query := `SELECT * FROM public.short_urls ORDER BY created_at DESC`
	err := DbSelect(query, &found)
	return found, err
}

// FindAllShortUrls 查找所有短链接
func FindAllShortUrlsByPage(page, size int) ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	if page < 0 {
		return found, nil
	}
	offset := (page - 1) * size
	query := `SELECT * FROM public.short_urls ORDER BY id DESC LIMIT $1 OFFSET $2`
	err := DbSelect(query, &found, size, offset)
	return found, err
}

// FindPagedShortUrls 分页查找短链接
func FindPagedShortUrls(url string, page int, size int) ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	offset := (page - 1) * size
	query := "SELECT * FROM public.short_urls u ORDER BY u.id DESC LIMIT $1 OFFSET $2"
	if !utils.EmptyString(url) {
		query := "SELECT * FROM public.short_urls u WHERE u.short_url = $1 ORDER BY u.id DESC LIMIT $2 OFFSET $3"
		var foundUrl core.ShortUrl
		err := DbGet(query, &foundUrl, url, size, offset)
		if !foundUrl.IsEmpty() {
			found = append(found, foundUrl)
		}
		return found, err
	}
	return found, DbSelect(query, &found, size, offset)
}

// InsertShortUrl 插入短链接
func InsertShortUrl(url core.ShortUrl) error {
	query := `INSERT INTO public.short_urls (short_url, dest_url, created_at, is_valid, memo,open_type)
	 VALUES(:short_url,:dest_url,:created_at,:is_valid,:memo,:open_type)`
	return DbNamedExec(query, url)
}

func splitLogsArray(array []core.AccessLog, size int) [][]core.AccessLog {
	var chunks [][]core.AccessLog
	for {
		if len(array) <= 0 {
			break
		}
		if len(array) < size {
			size = len(array)
		}
		chunks = append(chunks, array[0:size])
		array = array[size:]
	}
	return chunks
}
