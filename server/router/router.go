package router

import (
	"aigc-image/server/controller"
	"aigc-image/server/middleware"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (a *API) InitRouter() {
	r := a.Router
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	//允许跨域中间件
	r.Use(middleware.Cors())
	//gin日志中间件
	r.Use(gin.Logger())
	//gin恢复中间件
	r.Use(gin.Recovery())
	//引入requsesid中间件
	r.Use(requestid.New())
	//引入Prom监控中间件
	r.Use(ginprom.PromMiddleware(nil))
	//暴露Prom指标获取端口
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	//引用日志中间件
	//TODO 有fatal error:concurrent map read and map write的问题
	//r.Use(middleware.RequestLog())
	//设置GinMode
	//gin.SetMode(config.Confs.BaseConf.GinMode)

	dragonAPIGroup := r.Group("/dragon/api")
	{
		verifyGroup := dragonAPIGroup.Group("/algo")
		{
			verifyGroup.POST("/createImage", controller.CreateImage)
			verifyGroup.POST("/createProImage", controller.CreateProImage)
			verifyGroup.GET("/listRecords", controller.ListRecords)
		}
	}
}
