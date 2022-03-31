// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package service

import (
	"ohurlshortener/storage"
	"ohurlshortener/utils"
	"testing"
	"time"
)

func TestStoreAccessLog(t *testing.T) {
	init4Test(t)
	if err := StoreAccessLogs(); err != nil {
		t.Error(err)
	}
}

func TestNewAccessLog(t *testing.T) {
	init4Test(t)
	for i := 0; i < 100; i++ {
		if err := NewAccessLog("A1HeJzob", "192.168.2.1", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36", "fffff"); err != nil {
			t.Error(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func init4Test(t *testing.T) {
	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = storage.InitDatabaseService()
	if err != nil {
		t.Error(err)
		return
	}

	_, err = storage.InitRedisService()
	if err != nil {
		t.Error(err)
		return
	}
}
