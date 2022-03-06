package controller

import (
	"fmt"
	"net/http"
	"ohurlshortener/service"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
)

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
		"message": result,
		"result":  nil,
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
