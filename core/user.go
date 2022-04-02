package core

import "reflect"

type User struct {
	ID       int    `db:"id"`
	Account  string `db:"account"`
	Password string `db:"password"`
}

func (user User) IsEmpty() bool {
	return reflect.DeepEqual(user, User{})
}
