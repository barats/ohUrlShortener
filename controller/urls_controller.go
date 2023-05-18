// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package controller

import (
	"net/http"

	"ohurlshortener/core"
	"ohurlshortener/service"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
)

// ShortUrlDetail 重定向到目标地址
func ShortUrlDetail(c *gin.Context) {
	url := c.Param("url")
	if utils.EmptyString(url) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"code":    http.StatusNotFound,
			"message": "您访问的页面已失效",
			"label":   "Status Not Found",
		})
		return
	}

	memUrl, err := service.Search4ShortUrl(url)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "内部错误 - ohUrlShortener",
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"label":   "Error",
		})
		return
	}

	if utils.EmptyString(memUrl.DestUrl) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"code":    http.StatusNotFound,
			"message": "您访问的页面已失效",
			"label":   "Status Not Found",
		})
		return
	}

	ua := c.Request.UserAgent()
	switch ot := memUrl.OpenType; ot {
	case core.OpenInAndroid:
		if utils.IsAndroid(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInDingTalk:
		if utils.IsDingTalk(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInChrome:
		if utils.IsChrome(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInIPad:
		if utils.IsIPad(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInIPhone:
		if utils.IsIPhone(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInSafari:
		if utils.IsSafari(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInWeChat:
		if utils.IsWeChatUA(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInFirefox:
		if utils.IsFirefox(ua) {
			redirectSuccess(url, memUrl.DestUrl, c)
		} else {
			redirectFail(c)
		}
	case core.OpenInAll:
		redirectSuccess(url, memUrl.DestUrl, c)
	}
}

func redirectSuccess(shortUrl, destUrl string, ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, destUrl)
	go service.NewAccessLog(shortUrl, ctx.ClientIP(), ctx.Request.UserAgent(), ctx.Request.Referer())
}

func redirectFail(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "error.html", gin.H{
		"title":   "404 - ohUrlShortener",
		"code":    http.StatusNotFound,
		"message": "不支持的打开方式",
		"label":   "Status Not Found",
	})
}
