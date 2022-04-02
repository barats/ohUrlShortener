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
	"log"
	"net/http"
	"ohurlshortener/core"
	"ohurlshortener/service"
	"ohurlshortener/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminCookieValue(user core.User) (string, error) {
	var result string
	data, err := utils.Sha256Of(user.Account + "a=" + user.Password + "=e" + strconv.Itoa(user.ID))
	if err != nil {
		log.Println(err)
		return result, err
	}
	return utils.Base58Encode(data), nil
}

func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := c.Cookie("ohUrlShortenerAdmin")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
		}

		cookie, err := c.Cookie("ohUrlShortenerCookie")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
		}

		if len(user) <= 0 || len(cookie) <= 0 {
			c.Redirect(http.StatusFound, "/login")
		}

		found, err := service.GetUserByAccountFromRedis(user)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
		}

		if found.IsEmpty() {
			c.Redirect(http.StatusFound, "/login")
		}

		cValue, err := AdminCookieValue(found)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
		}

		if !strings.EqualFold(cValue, cookie) {
			c.Redirect(http.StatusFound, "/login")
		}

		c.Next()
	} //end of func
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
