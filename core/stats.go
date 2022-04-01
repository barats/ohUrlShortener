// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

type ShortUrlStats struct {
	ShortUrl               string `db:"short_url" json:"short_url"`
	TodayCount             int    `db:"today_count" json:"today_count"`
	YesterdayCount         int    `db:"yesterday_count" json:"yesterday_count"`
	Last7DaysCount         int    `db:"last_7_days_count" json:"last_7_days_count"`
	MonthlyCount           int    `db:"monthly_count" json:"monthly_count"`
	TotalCount             int    `db:"total_count" json:"total_count"`
	DistinctTodayCount     int    `db:"d_today_count" json:"d_today_count"`
	DistinctYesterdayCount int    `db:"d_yesterday_count" json:"d_yesterday_count"`
	DistinctLast7DaysCount int    `db:"d_last_7_days_count" json:"d_last_7_days_count"`
	DistinctMonthlyCount   int    `db:"d_monthly_count" json:"d_monthly_count"`
	DistinctTotalCount     int    `db:"d_total_count" json:"d_total_count"`
}

type Top25Url struct {
	ShortUrl
	ShortUrlStats
}

type UrlIpCountStats struct {
	ShortUrl
	ShortUrlStats
}
