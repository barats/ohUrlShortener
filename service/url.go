package service

import (
	"log"
	"ohurlshortener/core"
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"
	"time"
)

func ReloadUrls() (bool, error) {
	urls, err := db.FindAllShortUrls()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	err = redis.FlushDB()
	if err != nil {
		log.Println(err)
		return false, utils.RaiseError("内部错误，请联系管理员")
	}

	for _, url := range urls {
		if url.Valid {
			err := redis.Set4Ever(url.ShortUrl, url.DestUrl)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	} //end of for
	return true, nil
}

func Search4ShortUrl(shortUrl string) (string, error) {
	destUrl, err := redis.GetString(shortUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}
	return destUrl, nil
}

func GenerateShortUrl(destUrl string) (string, error) {
	shortUrl, err := core.GenerateShortLink(destUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	foundUrl, err := db.FindShortUrl(shortUrl)
	if err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	if !foundUrl.IsEmpty() {
		return shortUrl, nil
	}

	url := core.ShortUrl{
		DestUrl:   destUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now(),
		Valid:     true,
	}
	if err := db.InsertShortUrl(url); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	if err := redis.Set4Ever(shortUrl, destUrl); err != nil {
		log.Println(err)
		return "", utils.RaiseError("内部错误，请联系管理员")
	}

	return shortUrl, nil
}
