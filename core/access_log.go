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
	"time"
)

type AccessLog struct {
	ID         int64          `db:"id"`
	ShortUrl   string         `db:"short_url"`
	AccessTime time.Time      `db:"access_time"`
	Ip         sql.NullString `db:"ip"`
	UserAgent  sql.NullString `db:"user_agent"`
}
