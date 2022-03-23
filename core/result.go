// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

import (
	"net/http"
	"time"
)

type ResultJson struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Date    time.Time   `json:"date"`
}

func ResultJsonSuccess() ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Result:  nil,
		Date:    time.Now(),
	}
}

func ResultJsonSuccessWithData(data interface{}) ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Result:  data,
		Date:    time.Now(),
	}
}

func ResultJsonError(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusInternalServerError,
		Message: message,
		Result:  nil,
		Date:    time.Now(),
	}
}

func ResultJsonBadRequest(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusBadRequest,
		Message: message,
		Result:  nil,
		Date:    time.Now(),
	}
}
