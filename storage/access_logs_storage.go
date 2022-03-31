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

func FindAccessLogs(shortUrl string) ([]core.AccessLog, error) {
	found := []core.AccessLog{}
	query := "SELECT * FROM public.access_logs l WHERE l.short_url = $1 ORDER BY l.id DESC"
	err := DbSelect(query, &found, shortUrl)
	return found, err
}

func InsertAccessLogs(logs []core.AccessLog) error {
	if len(logs) <= 0 {
		return nil
	}
	query := `INSERT INTO public.access_logs (short_url, access_time, ip, user_agent) VALUES(:short_url,:access_time,:ip,:user_agent)`
	if len(logs) >= Max_Insert_Count {
		logsSlice := splitLogsArray(logs, Max_Insert_Count)
		for _, slice := range logsSlice {
			err := DbNamedExec(query, slice)
			if err != nil {
				return err
			}
		}
	}
	return DbNamedExec(query, logs)
}

func FindAccessLogsCount(url string) (int, error) {
	var rowCount int
	query := "SELECT count(l.id) FROM public.access_logs l"
	if !utils.EemptyString(url) {
		query = "SELECT count(l.id) FROM public.access_logs l WHERE l.short_url = $1"
		err := DbGet(query, &rowCount, url)
		return rowCount, err
	}
	return rowCount, DbGet(query, &rowCount)
}

func FindAllAccessLogs(url string, page int, size int) ([]core.AccessLog, error) {
	found := []core.AccessLog{}
	offset := (page - 1) * size
	query := "SELECT * FROM public.access_logs l ORDER BY l.id DESC LIMIT $1 OFFSET $2"
	if !utils.EemptyString(url) {
		query := "SELECT * FROM public.access_logs l WHERE l.short_url = $1 ORDER BY l.id DESC LIMIT $2 OFFSET $3"
		err := DbSelect(query, &found, url, size, offset)
		return found, err
	}
	return found, DbSelect(query, &found, size, offset)
}
