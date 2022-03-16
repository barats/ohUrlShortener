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

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed assets/* templates/*
	FS embed.FS

	config_file = "config.ini"

	group errgroup.Group
)

func init() {
	//Things MUST BE DONE before app starts
	_, err := utils.InitConfig(config_file)
	utils.ExitOnError("Config initialization failed.", err)

	_, err = redis.InitRedisService()
	utils.ExitOnError("Redis initialization failed.", err)

	_, err = db.InitDatabaseService()
	utils.ExitOnError("Database initialization failed.", err)

	_, err = service.ReloadUrls()
	utils.PrintOnError("Realod urls failed.", err)
}

func main() {

	router01, err := initializeRoute01()
	utils.ExitOnError("Router01 initialize failed.", err)

	router02, err := initializeRoute02()
	utils.ExitOnError("Router02 initialize failed.", err)

	serverWeb := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", utils.AppConfig.Port),
		Handler:      router01,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	serverAdmin := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", utils.AppConfig.AdminPort),
		Handler:      router02,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	group.Go(func() error {
		return serverWeb.ListenAndServe()
	})

	group.Go(func() error {
		return serverAdmin.ListenAndServe()
	})

	group.Go(func() error {
		return startTicker()
	})

	log.Printf("[ohUrlShortener] portal starts http://127.0.0.1:%d", utils.AppConfig.Port)
	log.Printf("[ohUrlShortener] admin starts http://127.0.0.1:%d", utils.AppConfig.AdminPort)

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
	router.Use(gin.Recovery(), utils.OhUrlShortenerLogFormat("Portal"))

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
	router.Use(gin.Recovery(), utils.OhUrlShortenerLogFormat("Admin"))

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

	router.GET("/login", controller.ShortUrlDetail)
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
	//sleep for 60s to make sure main process is running
	time.Sleep(60 * time.Second)

	//Clear redis cache every 3 minutes
	ticker := time.NewTicker(3 * time.Minute)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLogs(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
			return err
		}
		log.Println("[StoreAccessLog] Finish.")
	}
	return nil
}
