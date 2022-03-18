package db

import (
	"testing"
)

func TestGetAccessLogsCount(t *testing.T) {

	init4Test(t)

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase 1", args: args{""}, want: 10, wantErr: false},
		{name: "TestCase 2", args: args{"AvTkHZP7"}, want: 3, wantErr: false},
		{name: "TestCase 3", args: args{"abcd"}, want: 0, wantErr: false},
		{name: "TestCase 4", args: args{"AvtkhZP7"}, want: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAccessLogsCount(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccessLogsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAccessLogsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllAccessLogs(t *testing.T) {
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
		{name: "TestCase 1", args: args{url: "abc", page: 1, size: 2}, want: 0, wantErr: false},
		{name: "TestCase 2", args: args{url: "A1HeJzob", page: 1, size: 2}, want: 2, wantErr: false},
		{name: "TestCase 3", args: args{url: "A1HeJzob", page: 1, size: 5}, want: 5, wantErr: false},
		{name: "TestCase 4", args: args{url: "A1HeJzob", page: 2, size: 5}, want: 2, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAllAccessLogs(tt.args.url, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAccessLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GetAllAccessLogs() = %d, want %d", len(got), tt.want)
			}
		})
	}
}
