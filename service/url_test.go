package service

import (
	"log"
	"os"
	"testing"
)

func TestGenerateShortUrl(t *testing.T) {
	init4Test(t)
	if err := StoreAccessLogs(); err != nil {
		t.Error(err)
	}

	res, err := GenerateShortUrl("https://github.com/barats/ohurlshortener")
	if err != nil {
		t.Error(err)
	}

	log.SetOutput(os.Stdout)

	log.Println(res)
}
