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
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"ohurlshortener/controller"
	"ohurlshortener/service"
	"ohurlshortener/storage"
	"ohurlshortener/utils"

	"github.com/Masterminds/sprig"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const (
	WebReadTimeout  = 15 * time.Second
	WebWriteTimeout = 15 * time.Second

	// AccessLogCleanInterval 清理 Redis 中的访问日志的时间间隔
	AccessLogCleanInterval = 1 * time.Minute

	// Top25CalcInterval Top25 榜单计算间隔
	Top25CalcInterval = 5 * time.Minute

	// StatsSumCalcInterval 仪表盘页面中其他几个统计数据计算间隔
	StatsSumCalcInterval = 5 * time.Minute

	// StatsIpSumCalcInterval 全部访问日志分析统计的间隔
	StatsIpSumCalcInterval = 30 * time.Minute
)

var (
	//go:embed assets/* templates/*
	FS embed.FS

	group errgroup.Group

	cmdStart  string
	cmdConfig string
)

func main() {
	flag.StringVar(&cmdStart, "s", "", "starts ohUrlShortener service:  admin | portal ")
	flag.StringVar(&cmdConfig, "c", "config.ini", "config file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `ohUrlShortener version:%s
		Usage: ohurlshortener [-s admin|portal|<omit to start both>] [-c config_file_path]`, utils.Version)
		flag.PrintDefaults()
	}

	flag.Parse()

	initSettings()

	router01, err := initializeRoute01()
	utils.ExitOnError("Router01 initialize failed.", err)

	router02, err := initializeRoute02()
	utils.ExitOnError("Router02 initialize failed.", err)

	portal := &http.Server{
		Addr:         fmt.Sprintf(":%d", utils.AppConfig.Port),
		Handler:      router01,
		ReadTimeout:  WebReadTimeout,
		WriteTimeout: WebWriteTimeout,
	}

	admin := &http.Server{
		Addr:         fmt.Sprintf(":%d", utils.AppConfig.AdminPort),
		Handler:      router02,
		ReadTimeout:  WebReadTimeout,
		WriteTimeout: WebWriteTimeout,
	}

	if strings.EqualFold("admin", strings.TrimSpace(cmdStart)) {
		startAdmin(group, *admin)
	} else if strings.EqualFold("portal", strings.TrimSpace(cmdStart)) {
		startPortal(group, *portal)
	} else if utils.EmptyString(cmdStart) {
		startPortal(group, *portal)
		startAdmin(group, *admin)
	} else {
		flag.Usage()
	}

	err = group.Wait()
	utils.ExitOnError("Group failed,", err)
}

func initSettings() {
	// Things MUST BE DONE before app starts
	_, err := utils.InitConfig(cmdConfig)
	utils.ExitOnError("Config initialization failed.", err)

	rs, err := storage.InitRedisService()
	utils.ExitOnError("Redis initialization failed.", err)

	if strings.EqualFold("redis", strings.ToLower(utils.CaptchaConfig.Store)) {
		crs := storage.CaptchaRedisStore{KeyPrefix: "oh_captcha", Expiration: 1 * time.Minute, RedisService: rs}
		captcha.SetCustomStore(&crs)
	}

	_, err = storage.InitDatabaseService()
	storage.CallProcedureStatsTop25() // recalculate when ohUrlShortener starts
	storage.CallProcedureStatsSum()   // recalculate when ohUrlShortener starts
	utils.ExitOnError("Database initialization failed.", err)

	_, err = service.ReloadUrls()
	utils.PrintOnError("Reload urls failed.", err)

	err = service.ReloadUsers()
	utils.PrintOnError("Reload users failed.", err)
}

func startPortal(g errgroup.Group, server http.Server) {
	group.Go(func() error {
		log.Println("[StoreAccessLog] ticker starts to serve")
		return startTicker1()
	})

	group.Go(func() error {
		log.Printf("[ohUrlShortener] portal starts at http://localhost:%d", utils.AppConfig.Port)
		return server.ListenAndServe()
	})
}

func startAdmin(g errgroup.Group, server http.Server) {

	group.Go(func() error {
		log.Println("[Top25Urls] ticker starts to serve")
		return startTicker2()
	})

	group.Go(func() error {
		log.Println("[StatsIpSum] ticker starts to serve")
		return startTicker3()
	})

	group.Go(func() error {
		log.Println("[StatsSum] ticker starts to serve")
		return startTicker4()
	})

	group.Go(func() error {
		log.Printf("[ohUrlShortener] admin starts at http://localhost:%d", utils.AppConfig.AdminPort)
		return server.ListenAndServe()
	})
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
} // end of router01

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
	admin.GET("/dashboard", controller.DashboardPage)
	admin.GET("/urls", controller.UrlsPage)
	admin.GET("/stats", controller.StatsPage)
	admin.GET("/search_stats", controller.SearchStatsPage)
	admin.GET("/access_logs", controller.AccessLogsPage)
	admin.POST("/urls/generate", controller.GenerateShortUrl)
	admin.POST("/urls/state", controller.ChangeState)
	admin.POST("/urls/delete", controller.DeleteShortUrl)
	admin.POST("/access_logs_export", controller.AccessLogsExport)

	api := router.Group("/api", controller.APIAuthHandler())
	api.POST("/account", controller.APINewAdmin)
	api.PUT("/account/:account/update", controller.APIAdminUpdate)
	api.POST("/url", controller.APIGenShortUrl)
	api.GET("/url/:url", controller.APIUrlInfo)
	api.DELETE("/url/:url", controller.APIDeleteUrl)
	api.PUT("/url/:url/change_state", controller.APIUpdateUrl)
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404 - ohUrlShortener",
			"message": "您访问的页面已失效",
			"code":    http.StatusNotFound,
			"label":   "Error",
		})
	})
	return router, nil
} // end of router01

func startTicker1() error {
	redisTicker := time.NewTicker(AccessLogCleanInterval)
	for range redisTicker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLogs(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
	return nil
}

func startTicker2() error {
	top25Ticker := time.NewTicker(Top25CalcInterval)
	for range top25Ticker.C {
		log.Println("[Top25Urls Ticker] Start.")
		if err := storage.CallProcedureStatsTop25(); err != nil {
			log.Printf("Error while trying to calculate Top25Urls %s", err)
		}
		log.Println("[Top25Urls Ticker] Finish.")
	}
	return nil
}

func startTicker3() error {
	statsIpSumTicker := time.NewTicker(StatsIpSumCalcInterval)
	for range statsIpSumTicker.C {
		log.Println("[StatsIpSum Ticker] Start.")
		if err := storage.CallProcedureStatsIPSum(); err != nil {
			log.Printf("Error while trying to calculate StatsIpSum %s", err)
		}
		log.Println("[StatsIpSum Ticker] Finish.")
	}
	return nil
}

func startTicker4() error {
	statsSumTicker := time.NewTicker(StatsSumCalcInterval)
	for range statsSumTicker.C {
		log.Println("[StatsSum Ticker] Start.")
		if err := storage.CallProcedureStatsSum(); err != nil {
			log.Printf("Error while trying to calculate StatsSum %s", err)
		}
		log.Println("[StatsSum Ticker] Finish.")
	}
	return nil
}
