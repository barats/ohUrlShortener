package core

import (
	"reflect"
	"time"
)

// User 用户
type User struct {
	ID        int       `db:"id"`
	Account   string    `db:"account"`
	CreatedAt time.Time `db:"created_at"`
	Password  string    `db:"password"`
	Enabled   bool      `db:"is_enable" json:"is_enable"`
}

// IsEmpty 判断是否为空
func (user User) IsEmpty() bool {
	return reflect.DeepEqual(user, User{})
}
