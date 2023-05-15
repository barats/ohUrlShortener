// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package utils

import (
	"regexp"
)

func IsAndroid(ua string) bool {
	regex := regexp.MustCompile(`(?i)Android\/[\d.]+`)
	return regex.MatchString(ua)
}

func IsIPhone(ua string) bool {
	regex := regexp.MustCompile(`(?i)iPhone\/[\d.]+`)
	return regex.MatchString(ua)
}

func IsIPad(ua string) bool {
	regex := regexp.MustCompile(`(?i)iPad\/[\d.]+`)
	return regex.MatchString(ua)
}

func IsWeChatUA(ua string) bool {
	regex := regexp.MustCompile(`(?i)MicroMessenger\/[\d.]+`)
	return regex.MatchString(ua)
}

func IsDingTalk(ua string) bool {
	regex := regexp.MustCompile(`(?i)DingTalk\/[\d.]+`)
	return regex.MatchString(ua)
}

func IsSafari(ua string) bool {
	regex := regexp.MustCompile(`(?i)Version\/[\d.]+ Safari\/[\d.]+`)
	return regex.MatchString(ua)
}

func IsChrome(ua string) bool {
	regex := regexp.MustCompile(`(?i)Chrome\/[\d.]+ Safari`)
	return regex.MatchString(ua)
}

func IsFirefox(ua string) bool {
	regex := regexp.MustCompile(`(?i)Firefox\/[\d.]+`)
	return regex.MatchString(ua)
}
