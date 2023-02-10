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

// ResultJson 返回结果
type ResultJson struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Date    time.Time   `json:"date"`
}

// ResultJsonSuccess 返回成功结果
func ResultJsonSuccess() ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Status:  true,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonSuccessWithData 返回成功结果
func ResultJsonSuccessWithData(data interface{}) ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Result:  data,
		Status:  true,
		Date:    time.Now(),
	}
}

// ResultJsonError 返回错误结果
func ResultJsonError(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusInternalServerError,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonBadRequest 返回错误结果
func ResultJsonBadRequest(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusBadRequest,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonUnauthorized 返回错误结果
func ResultJsonUnauthorized(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusUnauthorized,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}
