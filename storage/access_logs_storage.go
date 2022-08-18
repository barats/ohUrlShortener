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

func DeleteAccessLogs(shortUrl string) error {
	query := `DELETE from public.access_logs WHERE short_url = :short_url`
	return DbNamedExec(query, shortUrl)
}

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

// FindAccessLogsCount
//
// Find Access Logs Count and Unique IP Count
//
// First return value is total_count, Second return value is unique_ip_count ip count
func FindAccessLogsCount(url string, start, end string) (int, int, error) {
	type LogsCount struct {
		TotalCount    int `db:"total_count"`
		UniqueIpCount int `db:"unique_ip_count"`
	}
	query := `SELECT count(l.id) as total_count, count(distinct(l.ip)) as unique_ip_count FROM public.access_logs l WHERE 1=1 `
	if !utils.EmptyString(url) {
		query += fmt.Sprintf(` AND l.short_url = '%s'`, url)
	}
	if !utils.EmptyString(start) {
		query += fmt.Sprintf(` AND l.access_time >= to_date('%s','YYYY-MM-DD')`, start)
	}
	if !utils.EmptyString(end) {
		query += fmt.Sprintf(` AND l.access_time < to_date('%s','YYYY-MM-DD')`, end)
	}
	var count LogsCount
	return count.TotalCount, count.UniqueIpCount, DbGet(query, &count)
}

func FindAllAccessLogs(url string, start, end string, page, size int) ([]core.AccessLog, error) {
	found := []core.AccessLog{}
	offset := (page - 1) * size
	query := `SELECT * FROM public.access_logs l WHERE 1=1 `
	if !utils.EmptyString(url) {
		query += fmt.Sprintf(` AND l.short_url = '%s'`, url)
	}
	if !utils.EmptyString(start) {
		query += fmt.Sprintf(` AND l.access_time >= to_date('%s','YYYY-MM-DD')`, start)
	}
	if !utils.EmptyString(end) {
		query += fmt.Sprintf(` AND l.access_time < to_date('%s','YYYY-MM-DD')`, end)
	}

	query += ` ORDER BY l.id DESC LIMIT $1 OFFSET $2`
	return found, DbSelect(query, &found, size, offset)
}

func FindAllAccessLogsByUrl(url string) ([]core.AccessLog, error) {
	found := []core.AccessLog{}
	query := "SELECT * FROM public.access_logs l ORDER BY l.id DESC"
	if !utils.EmptyString(url) {
		query := "SELECT * FROM public.access_logs l WHERE l.short_url = $1 ORDER BY l.id DESC"
		err := DbSelect(query, &found, url)
		return found, err
	}
	return found, DbSelect(query, &found)
}
