package storage

import (
	"database/sql"
	"ohurlshortener/core"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
)

func TestInsertShortUrls(t *testing.T) {
	init4Test(t)

	for i := 0; i < 10000; i++ {
		destUrl := faker.URL()
		shortUrl, _ := core.GenerateShortLink(destUrl)
		url := core.ShortUrl{DestUrl: destUrl, ShortUrl: shortUrl, CreatedAt: time.Now(), Valid: true, Memo: sql.NullString{String: destUrl, Valid: true}}
		err := InsertShortUrl(url)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteShortUrlWithAccessLogs(t *testing.T) {

	init4Test(t)

	url1 := core.ShortUrl{ShortUrl: "hello"}
	url2 := core.ShortUrl{ShortUrl: "hello"}
	url3 := core.ShortUrl{ShortUrl: "hello"}

	type args struct {
		shortUrl core.ShortUrl
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase1", args: args{url1}, wantErr: false},
		{name: "TestCase1", args: args{url2}, wantErr: false},
		{name: "TestCase1", args: args{url3}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteShortUrlWithAccessLogs(tt.args.shortUrl); (err != nil) != tt.wantErr {
				t.Errorf("DeleteShortUrlWithAccessLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
