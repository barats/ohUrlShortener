package main

import (
	"fmt"
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/utils"

	"github.com/gin-gonic/gin"
)

var config_file = "config.ini"

func init() {
	//Things MUST BE DONE before app starts
	_, err := utils.InitConfig(config_file)
	utils.ExitOnError("Config initialization failed.", err)

	_, err = redis.InitRedisService()
	utils.ExitOnError("Redis initialization failed.", err)

	_, err = db.InitDatabaseService()
	utils.ExitOnError("Database initialization failed.", err)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	err := r.Run(fmt.Sprintf(":%d", utils.AppConfig.Port))
	utils.ExitOnError("[ohUrlShortener] web service failed to start.", err)
}
