package controller

import (
	"fmt"
	"net/http"
	"ohurlshortener/core"
	"ohurlshortener/service"
	"ohurlshortener/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE_NUM  = 1
	DEFAULT_PAGE_SIZE = 20
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录 - ohUrlShortener",
	})
}

func DoLogin(c *gin.Context) {
	//TODO: Login logic
}

func DoLogout(c *gin.Context) {
	//TODO: Login logic
}

func DashbaordPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":       "仪表盘 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
	})
}

func GenerateShortUrl(c *gin.Context) {
	destUrl := c.PostForm("dest_url")
	memo := c.PostForm("memo")

	if utils.EemptyString(destUrl) {
		c.JSON(http.StatusBadRequest, core.ResultJsonError("目标链接不能为空！"))
		return
	}

	result, err := service.GenerateShortUrl(destUrl, memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.ResultJsonError(err.Error()))
		return
	}

	json := map[string]string{
		"short_url": fmt.Sprintf("%s%s", utils.AppConfig.UrlPrefix, result),
	}
	c.JSON(http.StatusOK, core.ResultJsonSuccessWithData(json))
}

func UrlsPage(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	strPage := c.DefaultQuery("page", strconv.Itoa(DEFAULT_PAGE_NUM))
	strSize := c.DefaultQuery("size", strconv.Itoa(DEFAULT_PAGE_SIZE))
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DEFAULT_PAGE_NUM
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DEFAULT_PAGE_SIZE
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

func AccessLogsPage(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	strPage := c.DefaultQuery("page", strconv.Itoa(DEFAULT_PAGE_NUM))
	strSize := c.DefaultQuery("size", strconv.Itoa(DEFAULT_PAGE_SIZE))
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DEFAULT_PAGE_NUM
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DEFAULT_PAGE_SIZE
	}
	logs, err := service.GetPagedAccessLogs(strings.TrimSpace(url), page, size)
	c.HTML(http.StatusOK, "access_logs.html", gin.H{
		"title":       "访问日志查询 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
		"error":       err,
		"logs":        logs,
		"page":        page,
		"size":        size,
		"prefix":      utils.AppConfig.UrlPrefix,
		"first_page":  page == 1,
		"last_page":   len(logs) < size,
		"url":         strings.TrimSpace(url),
	})
}
