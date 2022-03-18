package controller

import "github.com/gin-gonic/gin"

func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_url", c.Request.URL.Path)
		c.Next()
	}
}
