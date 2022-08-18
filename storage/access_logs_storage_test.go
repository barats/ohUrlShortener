package storage

import (
	"database/sql"
	"math/rand"
	"ohurlshortener/core"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
)

func randomTime() time.Time {
	min := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func TestNewAccessLog(t *testing.T) {
	init4Test(t)
	var logs []core.AccessLog
	for i := 0; i < 100000; i++ {
		log := core.AccessLog{ShortUrl: "AvTkHZP7",
			Ip:         sql.NullString{Valid: true, String: faker.IPv4()},
			AccessTime: randomTime(),
			UserAgent:  sql.NullString{Valid: true, String: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36"}}
		logs = append(logs, log)
	}
	t.Run("TestInsertAccessLogs", func(t *testing.T) {
		err := InsertAccessLogs(logs)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestFindAccessLogsCount(t *testing.T) {
	init4Test(t)

	type args struct {
		url   string
		start string
		end   string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{name: "TestCase 1", args: args{url: "", start: "", end: ""}, want: 18, want1: 7, wantErr: false},
		{name: "TestCase 2", args: args{url: "gkT39tb5", start: "", end: ""}, want: 6, want1: 4, wantErr: false},
		{name: "TestCase 3", args: args{url: "gkT39tb5", start: "2022-04-20", end: ""}, want: 2, want1: 1, wantErr: false},
		{name: "TestCase 3", args: args{url: "gkT39tb5", start: "2022-04-07", end: "2022-04-11"}, want: 3, want1: 3, wantErr: false},
		{name: "TestCase 1", args: args{url: "", start: "2022-04-01", end: ""}, want: 18, want1: 7, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := FindAccessLogsCount(tt.args.url, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAccessLogsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindAccessLogsCount() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindAccessLogsCount() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
