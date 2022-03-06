package controller

import (
	"log"
	"net/http"
	"ohurlshortener/service"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
)

func ShortUrlDetail(c *gin.Context) {
	url := c.Param("url")
	if utils.EemptyString(url) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"code":    http.StatusNotFound,
			"message": "您访问的页面已失效",
		})
		return
	}

	destUrl, err := service.Search4ShortUrl(url)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "错误 - ohUrlShortener",
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if utils.EemptyString(destUrl) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"code":    http.StatusNotFound,
			"message": "您访问的页面已失效",
		})
		return
	}

	err = service.NewAccessLog(url, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		log.Println(err)
	}

	c.Redirect(http.StatusFound, destUrl)
}
