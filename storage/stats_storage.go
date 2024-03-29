package storage

import (
	"ohurlshortener/core"
	"ohurlshortener/utils"
)

// GetUrlStats 获取短链接的访问量统计信息
func GetUrlStats(url string) (core.ShortUrlStats, error) {
	found := core.ShortUrlStats{}
	query := `select * from public.stats_ip_sum WHERE short_url = $1`
	err := DbGet(query, &found, url)
	return found, err
}

// GetUrlCount 获取短链接总数
func GetUrlCount() (int, error) {
	var (
		result int
		query  = `SELECT count(l.id) FROM public.short_urls l`
	)

	// query := `SELECT n_live_tup AS estimate_rows FROM pg_stat_all_tables WHERE relname = 'short_urls'`
	return result, DbGet(query, &result)
}

// GetSumOfUrlStats 获取所有短链接的访问量统计信息
func GetSumOfUrlStats() (core.ShortUrlStats, error) {
	query := `SELECT * FROM public.stats_sum`
	result := core.ShortUrlStats{}
	data := []core.StatsSum{}
	err := DbSelect(query, &data)
	if err != nil {
		return result, err
	}
	for _, v := range data {
		switch v.Key {
		case "today_count":
			result.TodayCount = v.Value
		case "yesterday_count":
			result.YesterdayCount = v.Value
		case "last_7_days_count":
			result.Last7DaysCount = v.Value
		case "monthly_count":
			result.MonthlyCount = v.Value
		case "d_today_count":
			result.DistinctTodayCount = v.Value
		case "d_yesterday_count":
			result.DistinctYesterdayCount = v.Value
		case "d_last_7_days_count":
			result.DistinctLast7DaysCount = v.Value
		case "d_monthly_count":
			result.DistinctMonthlyCount = v.Value
		}
	}
	return result, nil
}

// GetTop25 获取访问量前 25 的短链接
func GetTop25() ([]core.Top25Url, error) {
	query := `SELECT u.*,s.today_count AS today_count,s.d_today_count AS d_today_count FROM public.short_urls u , public.stats_top25 s WHERE u.short_url = s.short_url`
	found := []core.Top25Url{}
	return found, DbSelect(query, &found)
}

// FindPagedUrlIpCountStats 获取单个短链接的 IP 访问量统计信息
func FindPagedUrlIpCountStats(url string, page int, size int) ([]core.UrlIpCountStats, error) {
	found := []core.UrlIpCountStats{}
	offset := (page - 1) * size
	query := `SELECT s.*,u.id,u.dest_url,u.created_at,u.is_valid,u.memo FROM public.stats_ip_sum s , public.short_urls u WHERE u.short_url = s.short_url ORDER BY u.created_at DESC LIMIT $1 OFFSET $2`
	if !utils.EmptyString(url) {
		query := `SELECT s.*,u.id,u.dest_url,u.created_at,u.is_valid,u.memo
		FROM public.stats_ip_sum s , public.short_urls u WHERE u.short_url = s.short_url AND u.short_url = $1 ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`
		var foundUrl core.UrlIpCountStats
		err := DbGet(query, &foundUrl, url, size, offset)
		if !foundUrl.IsEmpty() {
			found = append(found, foundUrl)
		}
		return found, err
	}
	return found, DbSelect(query, &found, size, offset)
}

// CallProcedureStatsIPSum
// Call scheduled procedures to calculate stats result.
//
// Suggested time interval to call this procedure : 30 ~ 60 minutes
func CallProcedureStatsIPSum() error {
	query := `SELECT 1 AS r FROM p_stats_ip_sum()`
	var r int
	return DbGet(query, &r)
}

// CallProcedureStatsTop25
// Call scheduled procedures to calculate stats result.
//
// Suggested time interval to call this procedure 5 ~ 10 minutes
func CallProcedureStatsTop25() error {
	query := `SELECT 2 AS r FROM p_stats_top25()`
	var r int
	return DbGet(query, &r)
}

// CallProcedureStatsSum
// Call scheduled procedures to calculate stats result.
//
// Suggested time interval to call this procedure : 5 ~ 10 minutes
func CallProcedureStatsSum() error {
	query := `SELECT 3 AS r FROM p_stats_sum()`
	var r int
	return DbGet(query, &r)
}
