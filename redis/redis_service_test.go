package redis

import (
	"ohurlshortener/utils"
	"testing"

	oredis "github.com/go-redis/redis/v8"
)

func init4Test(t *testing.T) {

	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = InitRedisService()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSet(t *testing.T) {
	init4Test(t)
	err := Set4Ever("hello", "world")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetString(t *testing.T) {
	init4Test(t)
	rs, err := GetString("hello")
	if err == oredis.Nil {
		t.Errorf("GetString() found NOTHING.")
		return
	} else if err != nil {
		t.Error(err)
		return
	}
	if rs != "world" {
		t.Errorf("GetString() wanted %s, found %s", "world", rs)
		return
	}
}
