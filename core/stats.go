package core

type ShortUrlStats struct {
	ShortUrl               string `db:"short_url" json:"short_url"`
	TodayCount             int    `db:"today_count" json:"today_count"`
	YesterdayCount         int    `db:"yesterday_count" json:"yesterday_count"`
	MonthlyCount           int    `db:"monthly_count" json:"monthly_count"`
	TotalCount             int    `db:"total_count" json:"total_count"`
	DistinctTodayCount     int    `db:"d_today_count" json:"d_today_count"`
	DistinctYesterdayCount int    `db:"d_yesterday_count" json:"d_yesterday_count"`
	DistinctMonthlyCount   int    `db:"d_monthly_count" json:"d_monthly_count"`
	DistinctTotalCount     int    `db:"d_total_count" json:"d_total_count"`
}
