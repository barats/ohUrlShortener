// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"ohurlshortener/utils"
)

type OpenType int

const (
	OpenInAll OpenType = iota
	OpenInWeChat
	OpenInDingTalk
	OpenInIPhone
	OpenInAndroid
	OpenInIPad
	OpenInSafari
	OpenInChrome
	OpenInFirefox
)

// ShortUrl 短链接
type ShortUrl struct {
	ID        int64          `db:"id" json:"id"`
	ShortUrl  string         `db:"short_url" json:"short_url"`
	DestUrl   string         `db:"dest_url" json:"desc_url"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	Valid     bool           `db:"is_valid" json:"is_valid"`
	Memo      sql.NullString `db:"memo" json:"memo"`
	OpenType  OpenType       `db:"open_type" json:"open_type"`
}

type MemShortUrl struct {
	DestUrl  string
	OpenType OpenType
}

// IsEmpty 判断是否为空
func (url ShortUrl) IsEmpty() bool {
	return reflect.DeepEqual(url, ShortUrl{})
}

// GenerateShortLink 生成短链接
func GenerateShortLink(initialLink string) (string, error) {
	if utils.EmptyString(initialLink) {
		return "", fmt.Errorf("empty string")
	}
	urlHash, err := utils.Sha256Of(initialLink)
	if err != nil {
		return "", err
	}
	// number := new(big.Int).SetBytes(urlHash).Uint64()
	// str := utils.Base58Encode([]byte(fmt.Sprintf("%d", number)))
	str := utils.Base58Encode(urlHash)
	return str[:8], nil
}
