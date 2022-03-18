package controller

import (
	"fmt"
	"net/http"
	"ohurlshortener/service"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
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

func UrlsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "urls.html", gin.H{
		"title":       "短链接列表 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
	})
}

func AccessLogsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "access_logs.html", gin.H{
		"title":       "访问日志查询 - ohUrlShortener",
		"current_url": c.Request.URL.Path,
	})
}

func ShortUrlsStats(c *gin.Context) {
	url := c.Param("url")
	if utils.EemptyString(url) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "缺少参数 url",
			"result":  nil,
		})
		return
	}

	found, err := service.GetShortUrlStats(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"result":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"result":  found,
	})
}

func GetShortUrls(c *gin.Context) {
	urls, err := service.GetAllShortUrls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"result":  nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"result":  urls,
	})
}

func ReloadRedis(c *gin.Context) {
	result, err := service.ReloadUrls()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"result":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"result":  result,
	})
}

func GenerateShortUrl(c *gin.Context) {
	destUrl := c.PostForm("dest_url")
	if utils.EemptyString(destUrl) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "缺少参数 dest_url",
			"result":  nil,
		})
		return
	}

	shortUrl, err := service.GenerateShortUrl(destUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"result":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"result":  fmt.Sprintf("https://i.barats.cn/l/%s", shortUrl),
	})
}
