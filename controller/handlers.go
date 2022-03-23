// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_url", c.Request.URL.Path)
		c.Next()
	}
}

func WebLogFormatHandler(server string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if !strings.HasPrefix(param.Path, "/assets") {
			return fmt.Sprintf("[%s | %s] %s %s %d %s \t%s %s %s \n",
				server,
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		} //end of if
		return ""
	}) //end of formatter
} //end of func
