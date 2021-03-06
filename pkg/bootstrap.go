package egin

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/daodao97/egin/pkg/cache"
	"github.com/daodao97/egin/pkg/db"
	"github.com/daodao97/egin/pkg/middleware"
	"github.com/daodao97/egin/pkg/route"
	"github.com/daodao97/egin/pkg/utils"
)

type Bootstrap struct {
	HttpMiddlewares []func() gin.HandlerFunc
	RoutesMap       []route.RoutesMap
	RoutesGroup     []route.RoutesGroup
	engine          *gin.Engine
}

func (boot *Bootstrap) Start() {
	db.InitDb()
	cache.InitDb()
	gin.SetMode(utils.Config.Mode)
	// gin.DefaultWriter = ginLogger()
	boot.engine = gin.Default()
	boot.initValidator()
	boot.regMiddlewares()
	boot.regRoutes()
	err := boot.engine.Run(utils.Config.Address)
	if err != nil {
		return
	}
}

func (boot *Bootstrap) regMiddlewares() {
	for _, midFunc := range boot.HttpMiddlewares {
		boot.engine.Use(midFunc())
	}
	boot.engine.Use(gin.Recovery())
}

func (boot *Bootstrap) regRoutes() {
	route.RegRoutes(boot.engine, boot.RoutesMap)
	route.RegRouteGroup(boot.engine, boot.RoutesGroup)
	boot.engine.GET("/metrics", middleware.PromHandler(promhttp.Handler()))
	boot.engine.GET("/consul", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}

func (boot *Bootstrap) initValidator() {
	utils.InitValidator()
	for _, item := range boot.RoutesMap {
		for _, singleRoute := range item {
			for _, custom := range singleRoute.CustomValidateFuncs {
				utils.RegCustomValidateFunc(custom)
			}
		}
	}
	for _, item := range boot.RoutesGroup {
		for _, groupInfo := range item {
			for _, singleRoute := range groupInfo.RoutesMap {
				for _, custom := range singleRoute.CustomValidateFuncs {
					utils.RegCustomValidateFunc(custom)
				}
			}
		}
	}
}

func ginLogger() io.Writer {
	f, _ := os.Create("gin.log")
	return io.MultiWriter(f, os.Stdout)
}
