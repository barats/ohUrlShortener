package storage

import (
	"ohurlshortener/core"
	"ohurlshortener/utils"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func FindAllUsers() ([]core.User, error) {
	var found []core.User
	query := `SELECT * FROM public.users u`
	return found, DbSelect(query, &found)
}

func NewUser(account string, password string) error {
	query := `INSERT INTO public.users (account, "password") VALUES(:account,:password)`
	data, err := PasswordBase58Hash(password)
	if err != nil {
		return err
	}
	return DbNamedExec(query, core.User{Account: account, Password: data})
}

func UpdateUser(user core.User) error {
	query := `UPDATE public.users SET account = :account , "password" = :password WHERE id = :id`
	return DbNamedExec(query, user)
}

func FindUserByAccount(account string) (core.User, error) {
	var user core.User
	query := `SELECT * FROM public.users u WHERE lower(u.account) = $1`
	return user, DbGet(query, &user, strings.ToLower(account))
}

func PasswordBase58Hash(password string) (string, error) {
	data, err := utils.Sha256Of(password)
	if err != nil {
		return "", err
	}
	return base58.Encode(data), nil
}
