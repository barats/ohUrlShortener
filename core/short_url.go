package core

import "time"

type ShortUrl struct {
	ID        int64     `db:"id"`
	ShortUrl  string    `db:"short_url"`
	DestUrl   string    `db:"dest_url"`
	Sha       string    `db:"sha"`
	CreatedAt time.Time `db:"created_at"`
	Valid     bool      `db:"is_valid"`
}
