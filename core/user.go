package core

import "reflect"

// User 用户
type User struct {
	ID       int    `db:"id"`
	Account  string `db:"account"`
	Password string `db:"password"`
}

// IsEmpty 判断是否为空
func (user User) IsEmpty() bool {
	return reflect.DeepEqual(user, User{})
}
