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

func TestFindPagedShortUrls(t *testing.T) {
	init4Test(t)
	type args struct {
		url  string
		page int
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase 1", args: args{url: "", page: 1, size: 10}, want: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPagedShortUrls(tt.args.url, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPagedShortUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("FindPagedShortUrls() = %d, want %v", len(got), tt.want)
			}
		})
	}
}
