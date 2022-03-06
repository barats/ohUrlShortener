package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"ohurlshortener/controller"
	"ohurlshortener/db"
	"ohurlshortener/redis"
	"ohurlshortener/service"
	"ohurlshortener/utils"
	"time"

	"github.com/gin-gonic/gin"
)

var config_file = "config.ini"

//go:embed assets/* templates/*
var FS embed.FS

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
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	err := setupWebRoutes(r)
	utils.ExitOnError("Setup routes failed.", err)

	_, err = service.ReloadUrls()
	utils.PrintOnError("Realod urls failed.", err)

	go setupTicker()

	err = r.Run(fmt.Sprintf("127.0.0.1:%d", utils.AppConfig.Port))
	utils.ExitOnError("[ohUrlShortener] web service failed to start.", err)
}

func setupWebRoutes(router *gin.Engine) error {
	sub, err := fs.Sub(FS, "assets")
	if err != nil {
		return err
	}
	router.StaticFS("/assets", http.FS(sub))

	tmpl, err := template.New("").ParseFS(FS, "templates/*.html")
	if err != nil {
		return err
	}

	router.SetHTMLTemplate(tmpl)

	router.GET("/l/:url", controller.ShortUrlDetail)

	admin := router.Group("/admin")
	//TODO: 实现自己的用户验证逻辑，这里简单处理
	admin.Use(gin.BasicAuth(gin.Accounts{"ohUrlShortener": "helloworld"}))
	admin.GET("/shorturl", controller.GetShortUrls)
	admin.GET("/shorturl/:url/stats", controller.ShortUrlsStats)
	admin.POST("/shorturl", controller.GenerateShortUrl)

	admin.POST("/reload_redis", controller.ReloadRedis)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"message": "您访问的页面已失效",
		})
	})
	return nil
}

func setupTicker() {
	//sleep for 30s to make sure main process is gon
	time.Sleep(35 * time.Second)

	//Clear redis cache every 65 seconds
	ticker := time.NewTicker(65 * time.Second)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLog(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
}
