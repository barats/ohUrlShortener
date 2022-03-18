package service

import (
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"
	"testing"
	"time"
)

func TestStoreAccessLog(t *testing.T) {
	init4Test(t)
	if err := StoreAccessLogs(); err != nil {
		t.Error(err)
	}
}

func TestNewAccessLog(t *testing.T) {
	init4Test(t)
	for i := 0; i < 100; i++ {
		if err := NewAccessLog("A1HeJzob", "192.168.2.1", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36", "fffff"); err != nil {
			t.Error(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func init4Test(t *testing.T) {
	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.InitDatabaseService()
	if err != nil {
		t.Error(err)
		return
	}

	_, err = redis.InitRedisService()
	if err != nil {
		t.Error(err)
		return
	}
}
