// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package service

import (
	"log"
	"os"
	"testing"
)

func TestGenerateShortUrl(t *testing.T) {
	init4Test(t)
	if err := StoreAccessLogs(); err != nil {
		t.Error(err)
	}

	res, err := GenerateShortUrl("https://ww2222.ortener", "")
	if err != nil {
		t.Error(err)
	}

	log.SetOutput(os.Stdout)

	log.Println(res)
}
