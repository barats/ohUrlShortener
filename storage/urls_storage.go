// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package storage

import (
	"ohurlshortener/core"
	"ohurlshortener/utils"
)

var Max_Insert_Count = 5000

func UpdateShortUrl(shortUrl core.ShortUrl) error {
	query := `UPDATE public.short_urls SET short_url = :short_url, dest_url = :dest_url, is_valid = :is_valid, memo = :memo WHERE id = :id`
	return DbNamedExec(query, shortUrl)
}

func FindShortUrl(url string) (core.ShortUrl, error) {
	found := core.ShortUrl{}
	query := `SELECT * FROM public.short_urls WHERE short_url = $1`
	err := DbGet(query, &found, url)
	return found, err
}

func FindAllShortUrls() ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	query := `SELECT * FROM public.short_urls ORDER BY created_at DESC`
	err := DbSelect(query, &found)
	return found, err
}

func FindPagedShortUrls(url string, page int, size int) ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	offset := (page - 1) * size
	query := "SELECT * FROM public.short_urls u ORDER BY u.id DESC LIMIT $1 OFFSET $2"
	if !utils.EemptyString(url) {
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

func InsertShortUrl(url core.ShortUrl) error {
	query := `INSERT INTO public.short_urls (short_url, dest_url, created_at, is_valid, memo)
	 VALUES(:short_url,:dest_url,:created_at,:is_valid,:memo)`
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
