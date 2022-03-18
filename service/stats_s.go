package service

import (
	"log"
	"ohurlshortener/core"
	"ohurlshortener/db"
	"ohurlshortener/utils"
)

func GetShortUrlStats(url string) (core.ShortUrlStats, error) {
	found, err := db.GetUrlStats(url)
	if err != nil {
		log.Println(err)
		return found, utils.RaiseError("内部错误，请联系管理员")
	}
	return found, nil
}
