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
	"net/http"
	"strconv"
	"strings"
	"time"

	"ohurlshortener/core"
	"ohurlshortener/service"
	"ohurlshortener/utils"
	"ohurlshortener/utils/export"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPageNum  = 1
	DefaultPageSize = 20
)

// LoginPage 登录页面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录 - ohUrlShortener",
	})
}

// Users Page
func UsersPage(c *gin.Context) {
	strPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := c.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}

	found, err := service.GetPagedUsers(page, size)
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":       "用户管理 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
		"users":       found,
		"error":       err,
		"page":        page,
		"size":        size,
		"first_page":  page == 1,
		"last_page":   len(found) < size,
	})
}

// DoLogin 登录
func DoLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	captchaText := c.PostForm("captcha-text")
	captchaId := c.PostForm("captcha-id")

	if utils.EmptyString(account) || utils.EmptyString(password) || len(account) < 5 || len(password) < 8 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - ohUrlShortener",
			"error": "用户名或密码格式错误！",
		})
		return
	}

	if utils.EmptyString(captchaText) || utils.EmptyString(captchaId) || len(captchaText) < 6 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - ohUrlShortener",
			"error": "验证码格式错误!",
		})
		return
	}

	// 验证码有效性验证
	if !captcha.VerifyString(captchaId, captchaText) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - ohUrlShortener",
			"error": "验证码错误，请刷新页面再重新尝试！",
		})
		return
	}

	// 用户名密码有效性验证
	loginUser, err := service.Login(account, password)
	if err != nil || loginUser.IsEmpty() {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - ohUrlShortener",
			"error": err.Error(),
		})
		return
	}

	// Write Cookie to browser
	cValue, err := AdminCookieValue(loginUser)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "错误 - ohUrlShortener",
			"error": "内部错误，请联系管理员",
		})
		return
	}
	c.SetCookie("ohUrlShortenerAdmin", loginUser.Account, 3600, "/", "", false, true)
	c.SetCookie("ohUrlShortenerCookie", cValue, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// DoLogout 登出
func DoLogout(c *gin.Context) {
	c.SetCookie("ohUrlShortenerAdmin", "", -1, "/", "", false, true)
	c.SetCookie("ohUrlShortenerCookie", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// ServeCaptchaImage 生成验证码
func ServeCaptchaImage(c *gin.Context) {
	captcha.Server(200, 45).ServeHTTP(c.Writer, c.Request)
}

// RequestCaptchaImage 请求验证码
func RequestCaptchaImage(c *gin.Context) {
	imageId := captcha.New()
	c.JSON(http.StatusOK, core.ResultJsonSuccessWithData(imageId))
}

// ChangeState 修改状态
func ChangeState(c *gin.Context) {
	destUrl := c.PostForm("dest_url")
	enable := c.PostForm("enable")

	if utils.EmptyString(destUrl) {
		c.JSON(http.StatusBadRequest, core.ResultJsonError("目标链接不存在！"))
		return
	}

	destEnable, err := strconv.ParseBool(enable)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.ResultJsonError("参数不合法！"))
		return
	}

	result, er := service.ChangeState(destUrl, destEnable)
	if er != nil {
		c.JSON(http.StatusInternalServerError, core.ResultJsonError(er.Error()))
		return
	}

	c.JSON(http.StatusOK, core.ResultJsonSuccessWithData(result))
}

// DeleteShortUrl 删除短链接
func DeleteShortUrl(c *gin.Context) {
	url := c.PostForm("short_url")
	if utils.EmptyString(strings.TrimSpace(url)) {
		c.JSON(http.StatusBadRequest, core.ResultJsonError("目标链接不存在！"))
		return
	}

	err := service.DeleteUrlAndAccessLogs(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.ResultJsonError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, core.ResultJsonSuccess())
}

// GenerateShortUrl 生成短链接
func GenerateShortUrl(c *gin.Context) {
	destUrl := c.PostForm("dest_url")
	memo := c.PostForm("memo")
	strOpenType := c.PostForm("open_type")
	openType, err := strconv.Atoi(strOpenType)
	if err != nil {
		openType = int(core.OpenInAll)
	}

	if utils.EmptyString(destUrl) {
		c.JSON(http.StatusBadRequest, core.ResultJsonError("目标链接不能为空！"))
		return
	}

	result, err := service.GenerateShortUrl(destUrl, memo, openType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.ResultJsonError(err.Error()))
		return
	}

	json := map[string]string{
		"short_url": fmt.Sprintf("%s%s", utils.AppConfig.UrlPrefix, result),
	}
	c.JSON(http.StatusOK, core.ResultJsonSuccessWithData(json))
}

// StatsPage 统计页面
func StatsPage(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	strPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := c.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}
	urls, err := service.GetPagedUrlIpCountStats(strings.TrimSpace(url), page, size)
	c.HTML(http.StatusOK, "stats.html", gin.H{
		"title":       "数据统计 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
		"error":       err,
		"shortUrls":   urls,
		"page":        page,
		"size":        size,
		"prefix":      utils.AppConfig.UrlPrefix,
		"first_page":  page == 1,
		"last_page":   len(urls) < size,
		"url":         strings.TrimSpace(url),
	})
}

// SearchStatsPage 查询统计页面
func SearchStatsPage(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	strPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := c.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}
	urls, err := service.GetPagedUrlIpCountStats(strings.TrimSpace(url), page, size)
	c.HTML(http.StatusOK, "search_stats.html", gin.H{
		"title":       "查询统计 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
		"error":       err,
		"shortUrls":   urls,
		"page":        page,
		"size":        size,
		"prefix":      utils.AppConfig.UrlPrefix,
		"first_page":  page == 1,
		"last_page":   len(urls) < size,
		"url":         strings.TrimSpace(url),
	})
}

// UrlsPage 短链接列表页面
func UrlsPage(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	strPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := c.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}
	urls, err := service.GetPagesShortUrls(strings.TrimSpace(url), page, size)
	c.HTML(http.StatusOK, "urls.html", gin.H{
		"title":       "短链接列表 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
		"error":       err,
		"shortUrls":   urls,
		"page":        page,
		"size":        size,
		"prefix":      utils.AppConfig.UrlPrefix,
		"first_page":  page == 1,
		"last_page":   len(urls) < size,
		"url":         strings.TrimSpace(url),
	})
}

// AccessLogsPage 访问日志页面
func AccessLogsPage(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	strPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := c.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}

	totalCount, distinctIpCount, err := service.GetAccessLogsCount(strings.TrimSpace(url), start, end)
	logs, err := service.GetPagedAccessLogs(strings.TrimSpace(url), start, end, page, size)
	c.HTML(http.StatusOK, "access_logs.html", gin.H{
		"title":           "访问日志查询 - ohUrlShortener",
		"current_url":     c.Request.URL.Path,
		"error":           err,
		"logs":            logs,
		"page":            page,
		"size":            size,
		"prefix":          utils.AppConfig.UrlPrefix,
		"first_page":      page == 1,
		"last_page":       len(logs) < size,
		"url":             strings.TrimSpace(url),
		"total_count":     totalCount,
		"unique_ip_count": distinctIpCount,
		"start_date":      start,
		"end_date":        end,
	})
}

// AccessLogsExport 导出访问日志
func AccessLogsExport(c *gin.Context) {
	url := c.PostForm("url")
	logs, err := service.GetAllAccessLogs(strings.TrimSpace(url))

	if err != nil {
		c.HTML(http.StatusOK, "access_logs.html", gin.H{
			"title":       "访问日志查询 - ohUrlShortener",
			"current_url": c.Request.URL.Path,
			"prefix":      utils.AppConfig.UrlPrefix,
			"url":         strings.TrimSpace(url),
			"error":       err,
		})
		return
	}

	fileContent, err := export.AccessLogToExcel(logs)
	if err != nil {
		c.HTML(http.StatusOK, "access_logs.html", gin.H{
			"title":       "访问日志查询 - ohUrlShortener",
			"current_url": c.Request.URL.Path,
			"prefix":      utils.AppConfig.UrlPrefix,
			"url":         strings.TrimSpace(url),
			"error":       err,
		})
		return
	}

	attachmentName := "访问日志" + time.Now().Format("2006-01-02 150405") + ".xlsx"
	fileContentDisposition := "attachment;filename=\"" + attachmentName + "\""
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "pplication/octet-stream", fileContent)
}

// DashboardPage 仪表盘页面
func DashboardPage(c *gin.Context) {
	count, stats, err := service.GetSumOfUrlStats()
	if err != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":       "仪表盘 - ohUrlShortener",
			"current_url": c.Request.URL.Path,
			"error":       err,
		})
		return
	}

	top25, er := service.GetTop25Url()
	if er != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":       "仪表盘 - ohUrlShortener",
			"current_url": c.Request.URL.Path,
			"error":       er,
		})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":       "仪表盘 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
		"error":       err,
		"total_count": count,
		"prefix":      utils.AppConfig.UrlPrefix,
		"stats":       stats,
		"top25":       top25,
	})
}
