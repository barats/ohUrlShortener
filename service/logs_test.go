package service

import (
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"
	"testing"
)

func TestStoreAccessLog(t *testing.T) {
	init4Test(t)
	if err := StoreAccessLog(); err != nil {
		t.Error(err)
	}
}

func TestNewAccessLog(t *testing.T) {
	init4Test(t)
	if err := NewAccessLog("heh1e99999", "127.0.0.1", "asdfsdfas"); err != nil {
		t.Error(err)
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
