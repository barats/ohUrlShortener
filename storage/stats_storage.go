package storage

import (
	"ohurlshortener/core"
	"ohurlshortener/utils"
)

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

func FindPagedUrlIpCountStats(url string, page int, size int) ([]core.UrlIpCountStats, error) {
	found := []core.UrlIpCountStats{}
	offset := (page - 1) * size
	query := `SELECT s.*, u.id,u.dest_url,u.created_at,u.is_valid,u.memo
	FROM public.url_ip_count_stats s , public.short_urls u WHERE s.short_url = u.short_url ORDER BY u.created_at DESC LIMIT $1 OFFSET $2`
	if !utils.EemptyString(url) {
		query := `SELECT s.*, u.id,u.dest_url,u.created_at,u.is_valid,u.memo
		FROM public.url_ip_count_stats s , public.short_urls u WHERE s.short_url = u.short_url AND u.short_url = $1 ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`
		var foundUrl core.UrlIpCountStats
		err := DbGet(query, &foundUrl, url, size, offset)
		if !foundUrl.IsEmpty() {
			found = append(found, foundUrl)
		}
		return found, err
	}
	return found, DbSelect(query, &found, size, offset)
}
