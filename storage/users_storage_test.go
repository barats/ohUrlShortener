package storage

import (
	"ohurlshortener/utils"
	"testing"
)

func TestNewUser(t *testing.T) {
	init4Test(t)
	NewUser("ohUrlShortener", "-2aDzm=0(ln_9^1")
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
