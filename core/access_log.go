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
