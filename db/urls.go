package db

import "ohurlshortener/core"

var Max_Insert_Count = 5000

func FindShortUrl(url string) (core.ShortUrl, error) {
	found := core.ShortUrl{}
	query := `SELECT * FROM public.short_urls WHERE short_url = $1`
	err := Get(query, &found, url)
	return found, err
}

func FindAllShortUrls() ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	query := `SELECT * FROM public.short_urls ORDER BY created_at DESC`
	err := Select(query, &found)
	return found, err
}

func InsertShortUrl(url core.ShortUrl) error {
	query := `INSERT INTO public.short_urls
	(short_url, dest_url, created_at, is_valid)
	VALUES(:short_url,:dest_url,:created_at,:is_valid)`
	return NamedExec(query, url)
}

func FindAccessLogs(shortUrl string) ([]core.AccessLog, error) {
	found := []core.AccessLog{}
	query := "SELECT * FROM public.access_logs l WHERE l.short_url = $1 ORDER BY l.id DESC"
	err := Select(query, &found, shortUrl)
	return found, err
}

func InsertAccessLogs(logs []core.AccessLog) error {
	if len(logs) <= 0 {
		return nil
	}
	query := `INSERT INTO public.access_logs
	(short_url, access_time, ip, user_agent)
	VALUES(:short_url,:access_time,:ip,:user_agent)`
	if len(logs) >= Max_Insert_Count {
		logsSlice := splitLogsArray(logs, Max_Insert_Count)
		for _, slice := range logsSlice {
			err := NamedExec(query, slice)
			if err != nil {
				return err
			}
		}
	}
	return NamedExec(query, logs)
}

func GetUrlStats(url string) (core.ShortUrlStats, error) {
	found := core.ShortUrlStats{}
	query := `select * from public.url_ip_count_stats WHERE short_url = $1`
	err := Get(query, &found, url)
	return found, err
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
