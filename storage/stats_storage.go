package storage

import "ohurlshortener/core"

func GetUrlStats(url string) (core.ShortUrlStats, error) {
	found := core.ShortUrlStats{}
	query := `select * from public.url_ip_count_stats WHERE short_url = $1`
	err := DbGet(query, &found, url)
	return found, err
}

func GetUrlCount() (int, error) {
	var result int
	query := `SELECT count(l.id) FROM public.short_urls l`
	return result, DbGet(query, &result)
}

func GetSumOfUrlStats() (core.ShortUrlStats, error) {
	query := `SELECT * FROM public.sum_url_ip_count_stats`
	found := core.ShortUrlStats{}
	return found, DbGet(query, &found)
}

func GetTop25() ([]core.Top25Url, error) {
	query := `SELECT * FROM public.total_count_top25`
	found := []core.Top25Url{}
	return found, DbSelect(query, &found)
}
