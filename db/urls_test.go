package db

import (
	"database/sql"
	"log"
	"ohurlshortener/core"
	"testing"
	"time"
)

func TestInsertAccessLogs(t *testing.T) {
	init4Test(t)
	logs := []core.AccessLog{
		{AccessTime: time.Now(), ShortUrl: "http://gitee.com/barat"},
		{AccessTime: time.Now(), ShortUrl: "http://gitee.com/barat", UserAgent: sql.NullString{String: "asd user-agent", Valid: true}},
		{AccessTime: time.Now(), ShortUrl: "http://gitee.com/barat", UserAgent: sql.NullString{String: "a11sd user-agent", Valid: true},
			Ip: sql.NullString{String: "127.0.2", Valid: true}},
	}
	err := InsertAccessLogs(logs)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestFindAccessLogs(t *testing.T) {
	init4Test(t)

	logs, err := FindAccessLogs("https://gitee.com/barataaa")
	if err != nil {
		t.Error(err)
	}

	if len(logs) <= 0 {
		t.Error("could not find any record.")
	}
}

func TestFindShortUrl(t *testing.T) {
	init4Test(t)
	found, err := FindShortUrl("https://gitee.com/barat1")
	if err != nil {
		t.Error(err)
	}
	log.Printf("%v", found)
}
