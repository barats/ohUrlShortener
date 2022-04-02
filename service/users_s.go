package service

import (
	"encoding/json"
	"ohurlshortener/core"
	"ohurlshortener/storage"
	"ohurlshortener/utils"
	"strings"
)

const ADMIN_USER_PREFIX = "ohUrlShortenerAdmin#"
const ADMIN_COOKIE_PREFIX = "ohUrlShortenerCookie#"

func Login(account string, pasword string) (core.User, error) {

	var found core.User
	found, err := GetUserByAccountFromRedis(account)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员")
	}

	if found.IsEmpty() {
		return found, utils.RaiseError("用户名或密码错误")
	}

	res, err := storage.PasswordBase58Hash(pasword)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员")
	}

	if !strings.EqualFold(found.Password, res) {
		return found, utils.RaiseError("用户名或密码错误")
	}

	return found, nil
}

func ReloadUsers() error {
	users, err := storage.FindAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		jsonUser, _ := json.Marshal(user)
		er := storage.RedisSet4Ever(ADMIN_USER_PREFIX+user.Account, jsonUser)
		if er != nil {
			return er
		}
	}

	return nil
}

func GetUserByAccountFromRedis(account string) (core.User, error) {
	var found core.User
	foundUserStr, err := storage.RedisGetString(ADMIN_USER_PREFIX + account)
	if err != nil {
		return found, err
	}

	json.Unmarshal([]byte(foundUserStr), &found)

	return found, nil
}
