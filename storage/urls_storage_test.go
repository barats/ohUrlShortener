package storage

import (
	"ohurlshortener/core"
	"testing"
)

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
