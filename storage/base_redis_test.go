// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package storage

import (
	"ohurlshortener/utils"
	"reflect"
	"testing"
	"time"
)

func testInitRedisSettings(t *testing.T) {

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

func TestRedisSet(t *testing.T) {

	testInitRedisSettings(t)

	duration, _ := time.ParseDuration("30m")

	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test1", args: args{key: "hello", value: "test", ttl: duration}, wantErr: false},
		{name: "Test1", args: args{key: "world", value: "test1111", ttl: duration}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RedisSet(tt.args.key, tt.args.value, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("RedisSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisScan4Keys(t *testing.T) {

	want1 := []string{"hello"}
	want2 := []string{"world"}
	testInitRedisSettings(t)

	type args struct {
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Test1", args: args{prefix: "hel*"}, want: want1, wantErr: false},
		{name: "Test2", args: args{prefix: "wo*"}, want: want2, wantErr: false},
		{name: "Test2", args: args{prefix: "aa*"}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RedisScan4Keys(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisScan4Keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisScan4Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisGetString(t *testing.T) {

	testInitRedisSettings(t)
	want1 := "test"
	want2 := "test1111"

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Test1", args: args{key: "hello"}, want: want1, wantErr: false},
		{name: "Test2", args: args{key: "world"}, want: want2, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RedisGetString(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisGetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisGetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisFlushDB(t *testing.T) {
	testInitRedisSettings(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Test1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RedisFlushDB(); (err != nil) != tt.wantErr {
				t.Errorf("RedisFlushDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisDelete(t *testing.T) {

	testInitRedisSettings(t)

	type args struct {
		key []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test1", args: args{key: []string{"hello", "world"}}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RedisDelete(tt.args.key...); (err != nil) != tt.wantErr {
				t.Errorf("RedisDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
