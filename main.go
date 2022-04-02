// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"ohurlshortener/controller"
	"ohurlshortener/service"
	"ohurlshortener/storage"
	"ohurlshortener/utils"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const (
	CONFIG_FILE               = "config.ini"
	ACCESS_LOG_CLEAN_INTERVAL = 1 * time.Minute
	WEB_READ_TIMEOUT          = 10 * time.Second
	WEB_WRITE_TIMEOUT         = 10 * time.Second
)

var (
	//go:embed assets/* templates/*
	FS embed.FS

	group errgroup.Group
)

func init() {
	//Things MUST BE DONE before app starts
	_, err := utils.InitConfig(CONFIG_FILE)
	utils.ExitOnError("Config initialization failed.", err)

	_, err = storage.InitRedisService()
	utils.ExitOnError("Redis initialization failed.", err)

	_, err = storage.InitDatabaseService()
	utils.ExitOnError("Database initialization failed.", err)

	_, err = service.ReloadUrls()
	utils.PrintOnError("Realod urls failed.", err)

	err = service.ReloadUsers()
	utils.PrintOnError("Realod users failed.", err)
}

func main() {

	router01, err := initializeRoute01()
	utils.ExitOnError("Router01 initialize failed.", err)

	router02, err := initializeRoute02()
	utils.ExitOnError("Router02 initialize failed.", err)

	serverWeb := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", utils.AppConfig.Port),
		Handler:      router01,
		ReadTimeout:  WEB_READ_TIMEOUT,
		WriteTimeout: WEB_WRITE_TIMEOUT,
	}

	serverAdmin := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", utils.AppConfig.AdminPort),
		Handler:      router02,
		ReadTimeout:  WEB_READ_TIMEOUT,
		WriteTimeout: WEB_WRITE_TIMEOUT,
	}

	group.Go(func() error {
		log.Printf("[ohUrlShortener] portal starts at http://127.0.0.1:%d", utils.AppConfig.Port)
		return serverWeb.ListenAndServe()
	})

	group.Go(func() error {
		log.Printf("[ohUrlShortener] admin starts at http://127.0.0.1:%d", utils.AppConfig.AdminPort)
		return serverAdmin.ListenAndServe()
	})

	group.Go(func() error {
		log.Println("[ohUrlShortener] ticker starts to serve")
		return startTicker()
	})

	err = group.Wait()
	utils.ExitOnError("Group failed,", err)
}

func initializeRoute01() (http.Handler, error) {

	if utils.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery(), controller.WebLogFormatHandler("Portal"))

	sub, err := fs.Sub(FS, "assets")
	if err != nil {
		return nil, err
	}
	router.StaticFS("/assets", http.FS(sub))

	tmpl, err := template.New("").Funcs(sprig.FuncMap()).ParseFS(FS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	router.SetHTMLTemplate(tmpl)

	router.GET("/:url", controller.ShortUrlDetail)
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"message": "您访问的页面已失效",
			"code":    http.StatusNotFound,
			"label":   "Error",
		})
	})
	return router, nil
} //end of router01

func initializeRoute02() (http.Handler, error) {

	if utils.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery(), controller.WebLogFormatHandler("Admin"))

	sub, err := fs.Sub(FS, "assets")
	if err != nil {
		return nil, err
	}
	router.StaticFS("/assets", http.FS(sub))

	tmpl, err := template.New("").Funcs(sprig.FuncMap()).ParseFS(FS, "templates/**/*.html")
	if err != nil {
		return nil, err
	}

	router.SetHTMLTemplate(tmpl)

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	})
	router.GET("/login", controller.LoginPage)
	router.POST("/login", controller.DoLogin)
	router.GET("/captcha/:imageId", controller.ServeCaptchaImage)
	router.POST("/captcha", controller.RequestCaptchaImage)

	admin := router.Group("/admin", controller.AdminAuthHandler())
	admin.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/admin/dashboard")
	})
	admin.POST("/logout", controller.DoLogout)
	admin.GET("/dashboard", controller.DashbaordPage)
	admin.GET("/urls", controller.UrlsPage)
	admin.GET("/stats", controller.StatsPage)
	admin.GET("/access_logs", controller.AccessLogsPage)
	admin.POST("/urls/generate", controller.GenerateShortUrl)
	admin.POST("/urls/state", controller.ChangeState)
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"message": "您访问的页面已失效",
			"code":    http.StatusNotFound,
			"label":   "Error",
		})
	})
	return router, nil
} //end of router01

func startTicker() error {
	ticker := time.NewTicker(ACCESS_LOG_CLEAN_INTERVAL)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLogs(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
	return nil
}
