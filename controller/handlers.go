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
	"strconv"
	"strings"

	"ohurlshortener/core"
	"ohurlshortener/service"
	"ohurlshortener/storage"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
)

// APIAuthHandler Authorization for /api
func APIAuthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if utils.EmptyString(authHeader) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.ResultJsonUnauthorized("Authorization Header is empty"))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.ResultJsonUnauthorized("Invalid Authorization Header"))
			return
		}

		if fields[0] != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.ResultJsonUnauthorized("Unsupported Authorization Type"))
			return
		}

		token := fields[1]
		res, err := validateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ResultJsonError("Internal error"))
			return
		}

		if !res {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.ResultJsonUnauthorized("Authorization failed"))
			return
		}

		ctx.Next()
	}
}

// AdminCookieValue Generate cookie value for admin user
func AdminCookieValue(user core.User) (string, error) {
	var result string
	data, err := utils.Sha256Of(user.Account + "a=" + user.Password + "=e" + strconv.Itoa(user.ID))
	if err != nil {
		log.Println(err)
		return result, err
	}
	return utils.Base58Encode(data), nil
}

// AdminAuthHandler Authorization for /admin
func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		user, err := c.Cookie("ohUrlShortenerAdmin")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
		}

		cookie, err := c.Cookie("ohUrlShortenerCookie")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
			return
		}

		if len(user) <= 0 || len(cookie) <= 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
			return
		}

		found, err := service.GetUserByAccountFromRedis(user)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
			return
		}

		if found.IsEmpty() {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
			return
		}

		cValue, err := AdminCookieValue(found)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
			return
		}

		if !strings.EqualFold(cValue, cookie) {
			c.Redirect(http.StatusFound, "/login")
			c.Next()
			return
		}

		c.Next()
	} // end of func
}

// WebLogFormatHandler Customized log format for web
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
		} // end of if
		return ""
	}) // end of formatter
} // end of func

func validateToken(token string) (bool, error) {
	users, err := storage.FindAllUsers()
	if err != nil {
		return false, err
	}

	for _, u := range users {
		if u.Password == token {
			return true, nil
		}
	}

	return false, nil
}
