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
	"log"
	"testing"
	"time"

	"ohurlshortener/utils"

	"github.com/dchest/captcha"
)

func TestNewRedisStore(t *testing.T) {

	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}

	rs, err := InitRedisService()
	if err != nil {
		t.Error(err)
		return
	}

	s := CaptchaRedisStore{KeyPrefix: "sssTest", Expiration: 1 * time.Minute, RedisService: rs}

	captcha.SetCustomStore(&s)

	id := captcha.New()
	log.Println("Generated id is " + id)
}
