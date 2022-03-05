package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShortUrlDetail(c *gin.Context) {
	url := c.Param("url")
	if url == "hello" {
		c.Redirect(http.StatusFound, "https://github.com/barats")
	} else {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"message": "您访问页面已失效",
		})
	}
}
