// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

//
// Custom redis storage for captcha, According to https://github.com/dchest/captcha/blob/master/store.go
//
// An object implementing Store interface can be registered with SetCustomStore
// function to handle storage and retrieval of captcha ids and solutions for
// them, replacing the default memory store.
//
// It is the responsibility of an object to delete expired and used captchas
// when necessary (for example, the default memory store collects them in Set
// method after the certain amount of captchas has been stored.)

package storage

import (
	"fmt"
	"time"

	"ohurlshortener/utils"

	"github.com/dchest/captcha"
)

const (
	DefaultPrefixKey = "oh_captcha"
)

type CaptchaRedisStore struct {
	RedisService *RedisService
	Expiration   time.Duration
	KeyPrefix    string
}

func NewRedisStore(rs *RedisService, expiration time.Duration, prefixKey string) (captcha.Store, error) {
	s := new(CaptchaRedisStore)
	s.RedisService = rs
	s.Expiration = expiration
	s.KeyPrefix = prefixKey
	if utils.EmptyString(s.KeyPrefix) {
		s.KeyPrefix = DefaultPrefixKey
	}
	return s, nil
}

func (s *CaptchaRedisStore) Set(id string, digit []byte) {
	key := fmt.Sprintf("%s#%s", s.KeyPrefix, id)
	val, err := RedisGetString(key)
	if !utils.EmptyString(val) || err != nil {
		panic(fmt.Sprintf("RedisSet error for captcha key %s. %s", key, err))
	}
	RedisSet(key, digit, s.Expiration)
}

func (s *CaptchaRedisStore) Get(id string, clear bool) []byte {
	key := fmt.Sprintf("%s#%s", s.KeyPrefix, id)
	val, err := RedisGetString(key)
	if err != nil {
		return nil
	}

	if clear {
		RedisDelete(key)
	}

	return []byte(val)
}
