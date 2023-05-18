// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package utils

import (
	"testing"
)

func TestIsWeChatUA(t *testing.T) {

	ua1 := "mozilla/5.0 (linux; u; android 4.1.2; zh-cn; mi-one plus build/jzo54k) applewebkit/534.30 (khtml, like gecko) version/4.0 mobile safari/534.30 micromessenger/5.0.1.352"
	ua2 := "mozilla/5.0 (iphone; cpu iphone os 5_1_1 like mac os x) applewebkit/534.46 (khtml, like gecko) mobile/9b206 micromessenger/5.0"
	ua3 := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

	type args struct {
		ua string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test1", args: args{ua: ua1}, want: true},
		{name: "Test2", args: args{ua: ua2}, want: true},
		{name: "Test3", args: args{ua: ua3}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWeChatUA(tt.args.ua); got != tt.want {
				t.Errorf("IsWeChatUA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSafari(t *testing.T) {
	ua1 := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	ua2 := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Version/11 Safari/537.36"
	type args struct {
		ua string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test1", args: args{ua: ua1}, want: false},
		{name: "Test2", args: args{ua: ua2}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSafari(tt.args.ua); got != tt.want {
				t.Errorf("IsSafari() = %v, want %v", got, tt.want)
			}
		})
	}
}
