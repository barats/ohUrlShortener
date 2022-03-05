package db

import (
	"database/sql"
	"fmt"
	"log"
	"ohurlshortener/core"
	"ohurlshortener/utils"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func Test2(t *testing.T) {
	init4Test(t)

	query1 := `INSERT INTO public.access_logs
	(short_url, access_time, ip, user_agent)
	VALUES(:short_url,:access_time,:ip,:user_agent)`

	wanted1 := core.AccessLog{
		ShortUrl:   "https://gitee.com/barat",
		AccessTime: time.Now(),
		Ip: sql.NullString{
			String: "127.0.0.1",
			Valid:  true,
		},
		UserAgent: sql.NullString{
			String: "hello world",
			Valid:  true,
		},
	}

	wanted2 := core.AccessLog{
		ShortUrl:   "https://gitee.com/barat",
		AccessTime: time.Now(),
	}

	w := []core.AccessLog{wanted1, wanted2}
	err := NamedExec(query1, w)
	if err != nil {
		t.Error(err)
		return
	}

	query2 := `select * from public.access_logs`
	wanted3 := []core.AccessLog{}
	err = Select(query2, &wanted3)
	if err != nil {
		t.Error(err)
		return
	}
	if len(wanted3) <= 0 {
		t.Errorf("found 0 row but expected more.")
	}
}

func Test1(t *testing.T) {
	init4Test(t)

	wanted1 := core.ShortUrl{
		ShortUrl:  "https://gitee.com/barat",
		DestUrl:   "https://github.com/barats",
		CreatedAt: time.Now(),
		Sha:       fmt.Sprintf("a%d", time.Now().Unix()),
		Valid:     true,
	}

	wanted2 := core.ShortUrl{
		ShortUrl:  "https://gitee.com/barat",
		DestUrl:   "https://github.com/barats",
		CreatedAt: time.Now(),
		Sha:       fmt.Sprintf("b%d", time.Now().Unix()),
		Valid:     true,
	}
	query1 := `INSERT INTO public.short_urls
	(short_url, dest_url, sha, created_at, is_valid)
	VALUES(:short_url,:dest_url,:sha,:created_at,:is_valid)`
	err := NamedExec(query1, []core.ShortUrl{wanted1, wanted2})
	if err != nil {
		t.Error(err)
		return
	}

	query2 := "select * from public.short_urls where is_valid = true and sha = $1"
	found := core.ShortUrl{}
	err = Get(query2, &found, wanted2.Sha)
	if err != nil {
		t.Error(err)
		return
	}
	if found.Sha != wanted2.Sha {
		t.Errorf("wanted %v found %v", wanted2, found)
		return
	}

	query3 := "select * from public.short_urls where is_valid = true"
	found2 := []core.ShortUrl{}
	err = Select(query3, &found2)
	if err != nil {
		t.Error(err)
		return
	}
	log.Printf("found size: %d", len(found2))
}

func init4Test(t *testing.T) {
	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = InitDatabaseService()
	if err != nil {
		t.Error(err)
		return
	}
}
